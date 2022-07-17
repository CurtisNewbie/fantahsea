package data

import (
	"fantahsea/client"
	"fantahsea/config"
	. "fantahsea/err"
	"fantahsea/util"
	. "fantahsea/util"
	"fantahsea/web/dto"
	"fmt"
	"os"
	"time"

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
)

// ------------------------------- entity start

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
	IsDel      IS_DEL
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
func CreateGalleryImage(cmd *CreateGalleryImageCmd, user *User) error {
	imageNo := util.GenNo("IMG")

	if isCreated, e := isImgCreatedAlready(cmd.FileKey); isCreated || e != nil {
		if e != nil {
			return e
		}
		return NewWebErr("Image added already")
	}

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

		absPath := ResolveAbsFPath(cmd.GalleryNo, imageNo)
		log.Infof("Created GalleryImage record, downloading file from file-service to %s", absPath)

		// download the file from file-service
		if e := client.DownloadFile(cmd.FileKey, absPath); e != nil {
			return e
		}
		return nil
	})
	return te
}

// List gallery images
func ListGalleryImages(cmd *ListGalleryImagesCmd, user *User) (*ListGalleryImagesResp, error) {

	if hasAccess, err := HasAccessToGallery(user.UserNo, cmd.GalleryNo); err != nil || !hasAccess {
		if err != nil {
			return nil, err
		}
		return nil, NewWebErr("You are not allowed to access this gallery")
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

	return &ListGalleryImagesResp{ImageNos: imageNos, Paging: *dto.BuildResPage(&cmd.Paging, total)}, nil
}

/* Resolve download info for image */
func ResolveImageDInfo(imageNo string, user *User) (*ImageDInfo, error) {

	gi, err := findGalleryImage(imageNo)
	if err != nil {
		return nil, err
	}

	return &ImageDInfo{Name: gi.Name, Path: ResolveAbsFPath(gi.GalleryNo, gi.ImageNo)}, nil
}

// Resolve the absolute path to the image
func ResolveAbsFPath(galleryNo string, imageNo string) string {

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

	return dir + "/" + imageNo
}

/*
	-----------------------------------------------------------

	Helper methods

	-----------------------------------------------------------
*/

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
		return nil, NewWebErr("Image not found")
	}

	return &img, nil
}

/* Check whether the gallery image is created already */
func isImgCreatedAlready(fileKey string) (bool, error) {
	db := config.GetDB()

	var id int
	tx := db.Raw(`
		SELECT id FROM gallery_image
		WHERE file_key = ?
		AND is_del = 0
		`, fileKey).Scan(&id)

	if e := tx.Error; e != nil || tx.RowsAffected < 1 {
		return false, tx.Error
	}

	return true, nil
}
