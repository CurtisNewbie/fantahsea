package data

import (
	"fantahsea/config"
	"fantahsea/util"
	"log"
	"time"
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
	updateBy   string
	isDel      int8
}

// ------------------------------- entity end

// ------------------------------- table names

func (Gallery) TableName() string {
	return "gallery"
}

// -------------------------------

type CreateGalleryCmd struct {
	Name     string
	CreateBy string
}

// Create a new Gallery
func CreateGallery(cmd *CreateGalleryCmd) (*Gallery, error) {
	log.Printf("Creating gallery, cmd: %v\n", cmd)

	db := config.GetDB()
	gallery := Gallery{
		GalleryNo: util.GenNo("GAL"),
		Name:      cmd.Name,
		CreateBy:  cmd.CreateBy,
		updateBy:  cmd.CreateBy,
		isDel:     IS_DEL_N,
	}

	result := db.Create(&gallery)
	if result.Error != nil {
		return nil, result.Error
	}

	return &gallery, nil
}
