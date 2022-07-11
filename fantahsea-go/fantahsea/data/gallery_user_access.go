package data

import (
	"errors"
	"fantahsea/config"
	"time"

	"gorm.io/gorm"
)

// ------------------------------- entity start

// User's access to a Gallery
type GalleryUserAccess struct {
	ID         int64
	GalleryNo  string
	UserNo     string
	CreateTime time.Time
	CreateBy   string
	UpdateTime time.Time
	UpdateBy   string
	IsDel      int8
}

func (GalleryUserAccess) TableName() string {
	return "gallery_user_access"
}

// ------------------------------- entity end

/* Check if user has access to the gallery */
func HasAccessToGallery(userNo string, galleryNo string) bool {

	db := config.GetDB()

	// check if the user has access to the gallery
	var userAccess GalleryUserAccess

	tx := db.Where("gallery_no = ? and user_no = ?", galleryNo, userNo).First(&userAccess)
	if e := tx.Error; e != nil {

		// record not found
		if errors.Is(e, gorm.ErrRecordNotFound) {
			return false
		}
		return false
	}

	// galleryUserAccess may be logically deleted
	if IsDeleted(userAccess.IsDel) {
		return false
	}

	return true
}
