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
	UpdateBy   string
	IsDel      int8
}

func (GalleryUserAccess) TableName() string {
	return "gallery_user_access"
}

// ------------------------------- entity end
