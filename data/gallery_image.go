package data

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/curtisnewbie/fantahsea/client"
	"github.com/curtisnewbie/gocommon/config"
	"github.com/curtisnewbie/gocommon/dao"
	"github.com/curtisnewbie/gocommon/util"
	"github.com/curtisnewbie/gocommon/web/dto"
	"github.com/curtisnewbie/gocommon/weberr"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// GalleryImage.State
type ImgState string

const (
	// downloading from file-service
	DOWNLOADING ImgState = "DOWNLOADING"

	// processing
	PROCESSING ImgState = "PROCESSING"

	// ready to download
	READY ImgState = "READY"

	// 30mb is the maximum size for an image
	IMAGE_SIZE_THRESHOLD int64 = 30 * 1048576
)

var (
	imageNoCache = cache.New(10*time.Minute, 5*time.Minute)
	imageSuffix  = map[string]struct{}{"jpeg": {}, "jpg": {}, "gif": {}, "png": {}, "svg": {}, "bmp": {}}
)

// ------------------------------- entity start

type TransferGalleryImageInDirReq struct {
	// gallery no
	GalleryNo string `json:"galleryNo"`

	// file key of the directory
	FileKey string `json:"fileKey"`
}

// Image that belongs to a Gallery
type GalleryImage struct {
	ID         int64
	GalleryNo  string
	ImageNo    string
	Name       string
	FileKey    string
	CreateTime time.Time
	CreateBy   string
	UpdateTime time.Time
	UpdateBy   string
	IsDel      dao.IS_DEL
}

func (GalleryImage) TableName() string {
	return "gallery_image"
}

// ------------------------------- entity end

type ImageDInfo struct {
	Name string
	Path string
}

type ListGalleryImagesCmd struct {
	GalleryNo  string `json:"galleryNo"`
	dto.Paging `json:"pagingVo"`
}

type ListGalleryImagesResp struct {
	ImageNos []string   `json:"imageNos"`
	Paging   dto.Paging `json:"pagingVo"`
}

type CreateGalleryImageCmd struct {
	GalleryNo string `json:"galleryNo"`
	Name      string `json:"name"`
	FileKey   string `json:"fileKey"`
}

// Create a gallery image record
func CreateGalleryImage(cmd *CreateGalleryImageCmd, user *util.User) error {

	creator, err := FindGalleryCreator(cmd.GalleryNo)
	if err != nil {
		return err
	}

	if *creator != user.UserNo {
		return weberr.NewWebErr("You are not allowed to upload image to this gallery")
	}

	if isCreated, e := isImgCreatedAlready(cmd.GalleryNo, cmd.FileKey); isCreated || e != nil {
		if e != nil {
			return e
		}
		log.Infof("Image '%s' added already", cmd.Name)
		return nil
	}

	imageNo := util.GenNoL("IMG", 25)
	db := config.GetDB()
	te := db.Transaction(func(tx *gorm.DB) error {

		const sql string = `
			insert into gallery_image (gallery_no, image_no, name, file_key, create_by)
			values (?, ?, ?, ?, ?)
		`
		ct := tx.Exec(sql, cmd.GalleryNo, imageNo, cmd.Name, cmd.FileKey, user.Username)
		if ct.Error != nil {
			return ct.Error
		}

		absPath := ResolveAbsFPath(cmd.GalleryNo, imageNo, false)
		log.Infof("Created GalleryImage record, downloading file from file-service to %s", absPath)

		// download the file from file-service
		if e := client.DownloadFile(cmd.FileKey, absPath); e != nil {
			return e
		}

		// todo import a third-party golang library to compress image ?
		// compress the file using `convert` on linux
		// convert original.png -resize 256x original-thumbnail.png
		tnabs := absPath + "-thumbnail"
		out, err := exec.Command("convert", absPath, "-resize", "200x", tnabs).Output()
		log.Infof("Converted image, output: %s, absPath: %s", out, tnabs)
		if err != nil {
			return err
		}

		return nil
	})
	return te
}

// List gallery images
func ListGalleryImages(cmd *ListGalleryImagesCmd, user *util.User) (*ListGalleryImagesResp, error) {
	log.Printf("ListGalleryImages, cmd: %+v", cmd)

	if hasAccess, err := HasAccessToGallery(user.UserNo, cmd.GalleryNo); err != nil || !hasAccess {
		if err != nil {
			return nil, err
		}
		return nil, weberr.NewWebErr("You are not allowed to access this gallery")
	}

	const selectSql string = `
		select image_no from gallery_image 
		where gallery_no = ?
		and is_del = 0
		limit ?, ?
	`
	offset := dto.CalcOffset(&cmd.Paging)

	var imageNos []string
	tx := config.GetDB().Raw(selectSql, cmd.GalleryNo, offset, cmd.Paging.Limit).Scan(&imageNos)
	if tx.Error != nil {
		return nil, tx.Error
	}

	if imageNos == nil {
		imageNos = []string{}
	}

	fakeImageNos := []string{}
	for _, s := range imageNos {
		fakeImgNo := util.GenNoL("TKN", 25)
		imageNoCache.Set(fakeImgNo, s, cache.DefaultExpiration)
		fakeImageNos = append(fakeImageNos, fakeImgNo)
	}

	const countSql string = `
		select count(*) from gallery_image 
		where gallery_no = ?
		and is_del = 0
	`
	var total int
	tx = config.GetDB().Raw(countSql, cmd.GalleryNo).Scan(&total)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &ListGalleryImagesResp{ImageNos: fakeImageNos, Paging: *dto.BuildResPage(&cmd.Paging, total)}, nil
}

/* Resolve download info for image */
func ResolveImageDInfo(token string, thumbnail string) (*ImageDInfo, error) {

	imageNo, found := imageNoCache.Get(token)
	if !found {
		return nil, weberr.NewWebErr("You session has expired, please try again")
	}

	// log.Printf("Resolve Image DInfo, token: %s, imageNo: %s", token, imageNo)
	gi, e := findGalleryImage(imageNo.(string))
	if e != nil {
		return nil, e
	}

	return &ImageDInfo{Name: gi.Name, Path: ResolveAbsFPath(gi.GalleryNo, gi.ImageNo, strings.ToLower(thumbnail) == "true")}, nil
}

// Resolve the absolute path to the image
func ResolveAbsFPath(galleryNo string, imageNo string, thumbnail bool) string {
	basePath := config.GlobalConfig.FileConf.Base

	// convert to rune first
	runes := []rune(basePath)
	l := len(runes)
	if l < 1 {
		panic(fmt.Sprintf("unable to resolve image absolute path, base_path is illegal ('%x')", basePath))
	}

	if runes[l-1] != '/' {
		basePath += "/"
	}

	dir := basePath + galleryNo
	os.MkdirAll(dir, os.ModePerm)

	abs := dir + "/" + imageNo

	if thumbnail {
		abs = abs + "-thumbnail"
	}

	log.Printf("Resolved absolute path, galleryNo: %s, imageNo: %s, thumbnail: %t", galleryNo, imageNo, thumbnail)

	return abs
}

// Transfer images in dir
func TransferImagesInDir(req *TransferGalleryImageInDirReq, user *util.User) error {
	resp, e := client.GetFileInfo(req.FileKey)
	if e != nil {
		return e
	}

	fi := resp.Data

	// only the owner of the directory can do this, by default directory is only visible to the uploader
	if strconv.Itoa(fi.UploaderId) != user.UserId {
		return weberr.NewWebErr("Not permitted operation")
	}

	if fi.FileType != client.DIR {
		return weberr.NewWebErr("This is not a directory")
	}

	if fi.IsDeleted {
		return weberr.NewWebErr("Directory is already deleted")
	}

	go func(user *util.User, dirFileKey string, galleryNo string) {
		userNo := user.UserNo
		_, e := util.TimedLockRun("fantahsea:transfer:dir:"+userNo, 1*time.Second, func() any {
			start := time.Now()

			page := 1
			for true {
				resp, err := client.ListFilesInDir(dirFileKey, 100, page)
				if err != nil {
					log.Errorf("Failed to list files in dir, dir's fileKey: %s, error: %v", dirFileKey, err)
					break
				}
				if resp.Data == nil || len(resp.Data) < 1 {
					break
				}

				// starts fetching file one by one
				for i := 0; i < len(resp.Data); i++ {
					fk := resp.Data[i]
					fi, er := client.GetFileInfo(fk)
					if er != nil {
						log.Errorf("Failed to fetch file info while looping files in dir, fi's fileKey: %s, error: %v", fk, er)
						continue
					}

					if guessIsImage(fi.Data.Name, fi.Data.SizeInBytes) {
						if err := CreateGalleryImage(&CreateGalleryImageCmd{GalleryNo: galleryNo, Name: fi.Data.Name, FileKey: fi.Data.Uuid}, user); err != nil {
							log.Errorf("Failed to create gallery image, fi's fileKey: %s, error: %v", fk, err)
						}
					}
				}

				page += 1
			}

			log.Infof("Finished TransferImagesInDir, dir's fileKey: %s, took: %s", dirFileKey, time.Since(start))
			return nil
		})
		if e != nil && util.IsLockNotObtained(e) {
			log.Infof("Failed to obtain lock to transferImagesInDir, another goroutine may be transferring for current user, userNo: %s", userNo)
		}
	}(user, req.FileKey, req.GalleryNo)
	return nil
}

/*
	-----------------------------------------------------------

	Helper methods

	-----------------------------------------------------------
*/

// Guess whether a file is an image by its' name and size
func guessIsImage(name string, size int64) bool {
	if size > IMAGE_SIZE_THRESHOLD {
		return false
	}

	i := strings.LastIndex(name, ".")
	len := len([]rune(name))
	if i < 0 || i == len-1 {
		return false
	}

	suffix := name[i+1:]
	_, ok := imageSuffix[strings.ToLower(suffix)]
	return ok
}

/* Find gallery image */
func findGalleryImage(imageNo string) (*GalleryImage, error) {
	db := config.GetDB()

	var img GalleryImage
	tx := db.Raw(`
		SELECT * FROM gallery_image
		WHERE image_no = ?
		AND is_del = 0
		`, imageNo).Scan(&img)

	if e := tx.Error; e != nil {
		return nil, tx.Error
	}

	if tx.RowsAffected < 1 {
		log.Infof("Gallery Image not found, %s", imageNo)
		return nil, weberr.NewWebErr("Image not found")
	}

	return &img, nil
}

/* Check whether the gallery image is created already */
func isImgCreatedAlready(galleryNo string, fileKey string) (bool, error) {
	db := config.GetDB()

	var id int
	tx := db.Raw(`
		SELECT id FROM gallery_image
		WHERE gallery_no = ? 
		AND file_key = ?
		AND is_del = 0
		`, galleryNo, fileKey).Scan(&id)

	if e := tx.Error; e != nil || tx.RowsAffected < 1 {
		return false, tx.Error
	}

	return true, nil
}
