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
	Name     string
	CreateBy string
	UserNo   string
}

type UpdateGalleryCmd struct {
	GalleryNo string
	Name      string
	UpdateBy  string
	UserNo    string
}

// Create a new Gallery
func CreateGallery(cmd *CreateGalleryCmd) (*Gallery, error) {
	log.Printf("Creating gallery, cmd: %v\n", cmd)

	db := config.GetDB()
	gallery := Gallery{
		GalleryNo: util.GenNo("GAL"),
		Name:      cmd.Name,
		CreateBy:  cmd.CreateBy,
		UpdateBy:  cmd.CreateBy,
		IsDel:     IS_DEL_N,
	}

	result := db.Create(&gallery)
	if result.Error != nil {
		return nil, result.Error
	}

	return &gallery, nil
}

// Update a Gallery
func UpdateGallery(cmd *UpdateGalleryCmd) error {

	db := config.GetDB()
	glno := cmd.GalleryNo

	// check if the user has access to the gallery
	var userAccess GalleryUserAccess

	tx := db.Where("gallery_no = ? and user_no = ?", glno, cmd.UserNo).First(&userAccess)
	if e := tx.Error; e != nil {
		// record not found
		if errors.Is(e, gorm.ErrRecordNotFound) {
			return err.NewWebErr("You are not allowed to update this gallery")
		}
		return tx.Error
	}

	// galleryUserAccess may be logically deleted
	if IsDeleted(userAccess.IsDel) {
		return err.NewWebErr("You are not allowed to update this gallery")
	}

	tx = db.Where("gallery_no = ?", glno).Updates(Gallery{
		GalleryNo: cmd.GalleryNo,
		Name:      cmd.Name,
		UpdateBy:  cmd.UpdateBy,
	})

	if e := tx.Error; e != nil {
		return tx.Error
	}

	return nil
}
