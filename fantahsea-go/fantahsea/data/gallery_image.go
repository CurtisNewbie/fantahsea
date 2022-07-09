package data

import (
	"time"
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
	updateBy   string
	isDel      int8
}

// ------------------------------- entity end

// ------------------------------- table names

func (GalleryImage) TableName() string {
	return "gallery_image"
}

// -------------------------------
