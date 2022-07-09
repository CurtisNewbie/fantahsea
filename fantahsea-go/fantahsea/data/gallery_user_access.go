package data

import "time"

// ------------------------------- entity start

// User's access to a Gallery
type GalleryUserAccess struct {
	ID         int64
	GalleryNo  string
	UserNo     string
	CreateTime time.Time
	CreateBy   string
	UpdateTime time.Time
	updateBy   string
	isDel      int8
}

// ------------------------------- entity end

// ------------------------------- table names

func (GalleryUserAccess) TableName() string {
	return "gallery_user_access"
}

// -------------------------------
