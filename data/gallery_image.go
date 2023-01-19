package data

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/curtisnewbie/fantahsea/fclient"
	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/mysql"
	"github.com/curtisnewbie/gocommon/redis"
	"github.com/curtisnewbie/gocommon/server"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// GalleryImage.status
type ImgStatus string

const (
	NORMAL  ImgStatus = "NORMAL"
	DELETED ImgStatus = "DELETED"

	// 30mb is the maximum size for an image
	IMAGE_SIZE_THRESHOLD    int64 = 30 * 1048576
	DELETE_IMAGE_BATCH_SIZE int   = 30

	PROP_FILE_BASE = "file.base"
)

var (
	_imageNoCache *redis.RCache
	rwmu          sync.RWMutex

	imageSuffix = common.NewSet[string]()
)

func imageNoCache() *redis.RCache {
	rwmu.RLock()
	if _imageNoCache != nil {
		defer rwmu.RUnlock()
		return _imageNoCache
	}
	rwmu.RUnlock()

	rwmu.Lock()
	defer rwmu.Unlock()

	if _imageNoCache == nil {
		c := redis.NewRCache(1 * time.Minute)
		_imageNoCache = &c
	}
	return _imageNoCache
}

func init() {
	imageSuffix.Add("jpg")
	imageSuffix.Add("gif")
	imageSuffix.Add("png")
	imageSuffix.Add("svg")
	imageSuffix.Add("bmp")
	imageSuffix.Add("webp")
	imageSuffix.Add("apng")
	imageSuffix.Add("avif")
	common.SetDefProp(PROP_FILE_BASE, "files")
}

// ------------------------------- entity start

type TransferGalleryImageInDirReq struct {
	// gallery no
	GalleryNo string `json:"galleryNo" validation:"notEmpty"`

	// file key of the directory
	FileKey string `json:"fileKey" validation:"notEmpty"`
}

// Image that belongs to a Gallery
type GalleryImage struct {
	ID         int64
	GalleryNo  string
	ImageNo    string
	Name       string
	FileKey    string
	Status     ImgStatus
	CreateTime time.Time
	CreateBy   string
	UpdateTime time.Time
	UpdateBy   string
	IsDel      common.IS_DEL
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
	GalleryNo     string `json:"galleryNo" validation:"notEmpty"`
	common.Paging `json:"pagingVo"`
}

type ListGalleryImagesResp struct {
	ImageNos []string      `json:"imageNos"`
	Paging   common.Paging `json:"pagingVo"`
}

type CreateGalleryImageCmd struct {
	GalleryNo string `json:"galleryNo"`
	Name      string `json:"name"`
	FileKey   string `json:"fileKey"`
}

// Create a gallery image record
func CreateGalleryImage(ec server.ExecContext, cmd CreateGalleryImageCmd) error {
	user := ec.User
	creator, err := FindGalleryCreator(cmd.GalleryNo)
	if err != nil {
		return err
	}

	if *creator != user.UserNo {
		return common.NewWebErr("You are not allowed to upload image to this gallery")
	}

	if isCreated, e := isImgCreatedAlready(cmd.GalleryNo, cmd.FileKey); isCreated || e != nil {
		if e != nil {
			return e
		}
		ec.Log.Infof("Image '%s' added already", cmd.Name)
		return nil
	}

	imageNo := common.GenNoL("IMG", 25)
	db := mysql.GetMySql()
	te := db.Transaction(func(tx *gorm.DB) error {

		const sql string = `
			insert into gallery_image (gallery_no, image_no, name, file_key, create_by)
			values (?, ?, ?, ?, ?)
		`
		ct := tx.Exec(sql, cmd.GalleryNo, imageNo, cmd.Name, cmd.FileKey, user.Username)
		if ct.Error != nil {
			return ct.Error
		}

		absPath := ResolveAbsFPath(ec, cmd.GalleryNo, imageNo, false)
		ec.Log.Infof("Created GalleryImage record, downloading file from file-service to %s", absPath)

		// download the file from file-service
		if e := fclient.DownloadFile(ec.Ctx, cmd.FileKey, absPath); e != nil {
			return e
		}

		// todo import a third-party golang library to compress image ?
		// compress the file using `convert` on linux
		// convert original.png -resize 256x original-thumbnail.png
		tnabs := absPath + "-thumbnail"
		out, err := exec.Command("convert", absPath, "-resize", "200x", tnabs).Output()
		ec.Log.Infof("Converted image, output: %s, absPath: %s", out, tnabs)
		if err != nil {
			return err
		}

		return nil
	})
	return te
}

// List gallery images
func ListGalleryImages(cmd ListGalleryImagesCmd, ec server.ExecContext) (*ListGalleryImagesResp, error) {
	user := ec.User
	ec.Log.Infof("ListGalleryImages, cmd: %+v", cmd)

	if hasAccess, err := HasAccessToGallery(user.UserNo, cmd.GalleryNo); err != nil || !hasAccess {
		if err != nil {
			return nil, err
		}
		return nil, common.NewWebErr("You are not allowed to access this gallery")
	}

	const selectSql string = `
		select image_no from gallery_image 
		where gallery_no = ?
		and status = 'NORMAL'
		and is_del = 0
		order by id desc
		limit ?, ?
	`
	offset := common.CalcOffset(&cmd.Paging)

	var imageNos []string
	tx := mysql.GetMySql().Raw(selectSql, cmd.GalleryNo, offset, cmd.Paging.Limit).Scan(&imageNos)
	if tx.Error != nil {
		return nil, tx.Error
	}

	if imageNos == nil {
		imageNos = []string{}
	}

	fakeImageNos := []string{}
	for _, realImgNo := range imageNos {
		fakeImgNo := common.GenNoL("TKN", 25)
		e := imageNoCache().Put(fakeImgNo, realImgNo)
		if e != nil {
			return nil, e
		}
		fakeImageNos = append(fakeImageNos, fakeImgNo)
	}

	const countSql string = `
		select count(*) from gallery_image 
		where gallery_no = ?
		and status = 'NORMAL'
		and is_del = 0
	`
	var total int
	tx = mysql.GetMySql().Raw(countSql, cmd.GalleryNo).Scan(&total)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &ListGalleryImagesResp{ImageNos: fakeImageNos, Paging: *common.BuildResPage(&cmd.Paging, total)}, nil
}

/* Resolve download info for image */
func ResolveImageDInfo(token string, thumbnail string) (*ImageDInfo, error) {
	imageNo, e := imageNoCache().Get(token)
	if e != nil {
		return nil, e
	}

	if imageNo == "" {
		return nil, common.NewWebErr("You session has expired, please try again")
	}

	gi, e := findGalleryImage(imageNo)
	if e != nil {
		return nil, e
	}

	info := &ImageDInfo{
		Name: gi.Name,
		Path: ResolveAbsFPath(server.EmptyExecContext(), gi.GalleryNo, gi.ImageNo, strings.ToLower(thumbnail) == "true")}
	return info, nil
}

// Resolve the absolute path to the image
func ResolveAbsFPath(ec server.ExecContext, galleryNo string, imageNo string, thumbnail bool) string {
	basePath := common.GetPropStr("file.base")

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

	ec.Log.Printf("Resolved absolute path, galleryNo: %s, imageNo: %s, thumbnail: %t", galleryNo, imageNo, thumbnail)

	return abs
}

// Transfer images in dir
func TransferImagesInDir(cmd TransferGalleryImageInDirReq, ec server.ExecContext) error {
	user := ec.User
	resp, e := fclient.GetFileInfo(ec.Ctx, cmd.FileKey)
	if e != nil {
		return e
	}

	fi := resp.Data

	// only the owner of the directory can do this, by default directory is only visible to the uploader
	if strconv.Itoa(fi.UploaderId) != user.UserId {
		return common.NewWebErr("Not permitted operation")
	}

	if fi.FileType != fclient.DIR {
		return common.NewWebErr("This is not a directory")
	}

	if fi.IsDeleted {
		return common.NewWebErr("Directory is already deleted")
	}

	go func(ec server.ExecContext, user common.User, dirFileKey string, galleryNo string) {
		userNo := user.UserNo
		_, e := redis.TimedRLockRun("fantahsea:transfer:dir:"+userNo, 1*time.Second, func() (any, error) {
			start := time.Now()

			page := 1
			for {
				resp, err := fclient.ListFilesInDir(ec.Ctx, dirFileKey, 100, page)
				if err != nil {
					ec.Log.Errorf("Failed to list files in dir, dir's fileKey: %s, error: %v", dirFileKey, err)
					break
				}
				if resp.Data == nil || len(resp.Data) < 1 {
					break
				}

				// starts fetching file one by one
				for i := 0; i < len(resp.Data); i++ {
					fk := resp.Data[i]
					fi, er := fclient.GetFileInfo(ec.Ctx, fk)
					if er != nil {
						ec.Log.Errorf("Failed to fetch file info while looping files in dir, fi's fileKey: %s, error: %v", fk, er)
						continue
					}

					if guessIsImage(fi.Data.Name, fi.Data.SizeInBytes) {
						cmd := CreateGalleryImageCmd{GalleryNo: galleryNo, Name: fi.Data.Name, FileKey: fi.Data.Uuid}
						if err := CreateGalleryImage(ec, cmd); err != nil {
							ec.Log.Errorf("Failed to create gallery image, fi's fileKey: %s, error: %v", fk, err)
						}
					}
				}

				page += 1
			}

			ec.Log.Infof("Finished TransferImagesInDir, dir's fileKey: %s, took: %s", dirFileKey, time.Since(start))
			return nil, nil
		})
		if e != nil {
			ec.Log.Infof("Failed to transferImagesInDir, userNo: %s, dirFileKey: %s, galleryNo: %s, err: %v", userNo, dirFileKey, galleryNo, e)
		}
	}(ec, user, cmd.FileKey, cmd.GalleryNo)
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
	return imageSuffix.Has(strings.ToLower(suffix))
}

// Find gallery image
func findGalleryImage(imageNo string) (*GalleryImage, error) {
	db := mysql.GetMySql()

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
		logrus.Infof("Gallery Image not found, %s", imageNo)
		return nil, common.NewWebErr("Image not found")
	}

	return &img, nil
}

// check whether the gallery image is created already
//
// return isImgCreated, error
func isImgCreatedAlready(galleryNo string, fileKey string) (bool, error) {
	db := mysql.GetMySql()

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

// mark image as deleted
func markImageAsDeleted(imageNo string) error {
	tx := mysql.GetMySql().Exec(`
		update gallery_image
		set status = ?
		where image_no = ?
		`, DELETED, imageNo)

	if e := tx.Error; e != nil {
		return tx.Error
	}

	return nil
}

// find normal images of gallery
//
// return *[]imageNos, error
func findNormalImagesOfGallery(galleryNo string, limit int) (*[]string, error) {
	var imageNos []string
	tx := mysql.GetMySql().Raw(`
		select gi.image_no from gallery_image gi
		where gallery_no = ?
		and gi.status = 'NORMAL'
		and gi.is_del = 0
		limit ?
		`, galleryNo, limit).Scan(&imageNos)

	if e := tx.Error; e != nil || tx.RowsAffected < 1 {
		return nil, tx.Error
	}
	return &imageNos, nil
}

// find one deleted gallery that needs clean-up, i.e., gallery that still has images not deleted
//
// return *galleryNo, error
func findOneGalleryNeedsCleanup() (*string, error) {
	var gno string
	tx := mysql.GetMySql().Raw(`
		select g.gallery_no from gallery g
		where g.is_del = 1
		and exists (
			select * from gallery_image gi 
			where gi.gallery_no = g.gallery_no and gi.status = 'NORMAL'
		) 
		limit 1
		`).Scan(&gno)

	if e := tx.Error; e != nil || tx.RowsAffected < 1 {
		return nil, tx.Error
	}
	return &gno, nil
}

// clean up deleted galleries
func CleanUpDeletedGallery() {
	galleryNo, e := findOneGalleryNeedsCleanup()
	if e != nil {
		logrus.Errorf("Failed to find gallery that needs cleanup, err: %v", e)
		return
	}

	if galleryNo == nil {
		logrus.Infof("Found no gallery that needs clean-up")
		return
	}

	ec := server.EmptyExecContext()

	logrus.Infof("Found deleted gallery that needs clean-up, galleryNo: %s", *galleryNo)
	for {
		imageNos, err := findNormalImagesOfGallery(*galleryNo, DELETE_IMAGE_BATCH_SIZE)
		if err != nil {
			logrus.Errorf("Failed to find normal images of gallery, err: %v", err)
			return
		}
		if imageNos == nil || len(*imageNos) < 1 {
			break
		}

		for _, n := range *imageNos {
			imgNo := n

			img := ResolveAbsFPath(ec, *galleryNo, imgNo, false)
			if e := tryDeleteFile(img); e != nil {
				logrus.Errorf("Failed to delete file: %s, galleryNo: %s, err: %v", img, *galleryNo, e)
				return
			}

			thumbnail := ResolveAbsFPath(ec, *galleryNo, imgNo, true)
			if e := tryDeleteFile(thumbnail); e != nil {
				logrus.Errorf("Failed to delete file: %s, galleryNo: %s, err: %v", img, *galleryNo, e)
				return
			}

			if err := markImageAsDeleted(imgNo); err != nil {
				logrus.Errorf("Failed to mark image as deleted, %s, e: %v", imgNo, err)
			} else {
				logrus.Infof("Image deleted, %s", imgNo)
			}
		}
	}
	logrus.Infof("Finished deleting images of gallery, galleryNo: %s", *galleryNo)
}

// try to delete the file using os.Remove, if the file is deleted or not found, nil is returned, else the error
func tryDeleteFile(path string) error {
	if e := os.Remove(path); e != nil {
		if errors.Is(e, fs.ErrNotExist) {
			logrus.Infof("File is not found or already deleted, path: %s", path)
			return nil // the file is deleted already
		}
		return e
	}
	return nil
}
