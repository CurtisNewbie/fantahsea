package data

import (
	"errors"
	"fantahsea/config"
	"fantahsea/err"
	"fantahsea/util"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// ------------------------------- entity start

// Gallery
type Gallery struct {
	ID         int64
	GalleryNo  string
	UserNo     string
	Name       string
	CreateTime time.Time
	CreateBy   string
	UpdateTime time.Time
	UpdateBy   string
	IsDel      int8
}

func (Gallery) TableName() string {
	return "gallery"

}

// ------------------------------- entity end

type CreateGalleryCmd struct {
	Name string
}

type UpdateGalleryCmd struct {
	GalleryNo string
	Name      string
}

// Create a new Gallery
func CreateGallery(cmd *CreateGalleryCmd, user *util.User) (*Gallery, error) {
	log.Printf("Creating gallery, cmd: %v, user: %v\n", cmd, user)

	db := config.GetDB()
	gallery := &Gallery{
		GalleryNo: util.GenNo("GAL"),
		Name:      cmd.Name,
		UserNo:    user.UserNo,
		CreateBy:  user.Username,
		UpdateBy:  user.Username,
		IsDel:     IS_DEL_N,
	}

	result := db.Create(gallery)
	if result.Error != nil {
		return nil, result.Error
	}

	return gallery, nil
}

// Update a Gallery
func UpdateGallery(cmd *UpdateGalleryCmd, user *util.User) error {

	db := config.GetDB()
	galleryNo := cmd.GalleryNo

	gallery, e := FindGallery(galleryNo)
	if e != nil {
		return e
	}

	// only owner can update the gallery
	if user.UserNo != gallery.UserNo {
		return err.NewWebErr("You are not allowed to update this gallery")
	}

	tx := db.Where("gallery_no = ?", galleryNo).Updates(Gallery{
		GalleryNo: cmd.GalleryNo,
		Name:      cmd.Name,
		UpdateBy:  user.Username,
	})

	if e := tx.Error; e != nil {
		log.Warnf("Failed to update gallery, gallery_no: %v, e: %v\n", galleryNo, tx.Error)
		return err.NewWebErr("Failed to update gallery, please try again later")
	}

	return nil
}

/** Find Gallery by gallery_no */
func FindGallery(galleryNo string) (*Gallery, error) {

	db := config.GetDB()
	var gallery *Gallery
	tx := db.Where("gallery_no = ?", galleryNo).First(gallery)

	if e := tx.Error; e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			return nil, err.NewWebErr("Gallery doesn't exist")
		}
		return nil, tx.Error
	}
	return gallery, nil
}
