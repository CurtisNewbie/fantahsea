package data

import (
	"time"

	"github.com/curtisnewbie/gocommon/config"
	"github.com/curtisnewbie/gocommon/dao"
	"github.com/curtisnewbie/gocommon/util"
	"github.com/curtisnewbie/gocommon/web/dto"
	"github.com/curtisnewbie/gocommon/weberr"

	log "github.com/sirupsen/logrus"
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
	IsDel      dao.IS_DEL
}

func (Gallery) TableName() string {
	return "gallery"

}

// ------------------------------- entity end

type CreateGalleryCmd struct {
	Name string `json:"name"`
}

type UpdateGalleryCmd struct {
	GalleryNo string `json:"galleryNo"`
	Name      string `json:"name"`
}

type ListGalleriesResp struct {
	Paging    *dto.Paging `json:"pagingVo"`
	Galleries *[]VGallery `json:"galleries"`
}

type ListGalleriesCmd struct {
	Paging *dto.Paging `json:"pagingVo"`
}

type DeleteGalleryCmd struct {
	GalleryNo string `json:"galleryNo"`
}

type PermitGalleryAccessCmd struct {
	GalleryNo string `json:"galleryNo"`
	UserNo    string `json:"userNo"`
}

type VGallery struct {
	ID         int64     `json:"id"`
	GalleryNo  string    `json:"galleryNo"`
	UserNo     string    `json:"userNo"`
	Name       string    `json:"name"`
	CreateTime dto.WTime `json:"createTime"`
	CreateBy   string    `json:"createBy"`
	UpdateTime dto.WTime `json:"updateTime"`
	UpdateBy   string    `json:"updateBy"`
}

/* List Galleries */
func ListGalleries(cmd *ListGalleriesCmd, user *util.User) (*ListGalleriesResp, error) {
	paging := cmd.Paging

	const selectSql string = `
		SELECT g.* from gallery g 
		WHERE g.user_no = ? 
		AND g.is_del = 0 
		OR EXISTS (SELECT * FROM gallery_user_access ga WHERE ga.gallery_no = g.gallery_no AND ga.user_no = ?)
		LIMIT ?, ?
	`
	db := config.GetDB()
	var galleries []VGallery

	offset := dto.CalcOffset(paging)
	tx := db.Raw(selectSql, user.UserNo, user.UserNo, offset, paging.Limit).Scan(&galleries)

	if e := tx.Error; e != nil {
		return nil, e
	}

	const countSql string = `
		SELECT count(*) from gallery g 
		WHERE g.user_no = ? 
		AND g.is_del = 0 
		OR EXISTS (SELECT * FROM gallery_user_access ga WHERE ga.gallery_no = g.gallery_no AND ga.user_no = ?)
	`
	var total int
	tx = db.Raw(countSql, user.UserNo, user.UserNo).Scan(&total)

	if e := tx.Error; e != nil {
		return nil, e
	}

	if galleries == nil {
		galleries = []VGallery{}
	}

	return &ListGalleriesResp{Galleries: &galleries, Paging: dto.BuildResPage(paging, total)}, nil
}

// Create a new Gallery
func CreateGallery(cmd *CreateGalleryCmd, user *util.User) (*Gallery, error) {
	log.Printf("Creating gallery, cmd: %v, user: %v", cmd, user)

	// Guest is not allowed to create gallery
	if util.IsGuest(user) {
		return nil, weberr.NewWebErr("Guest is not allowed to create gallery")
	}

	galleryNo := util.GenNo("GAL")

	db := config.GetDB().Begin()
	gallery := &Gallery{
		GalleryNo: galleryNo,
		Name:      cmd.Name,
		UserNo:    user.UserNo,
		CreateBy:  user.Username,
		UpdateBy:  user.Username,
		IsDel:     dao.IS_DEL_N,
	}

	result := db.Omit("CreateTime", "UpdateTime").Create(gallery)
	if result.Error != nil {
		db.Rollback()
		return nil, result.Error
	}

	tx := db.Exec(`INSERT INTO gallery_user_access (gallery_no, user_no, create_by) VALUES (?, ?, ?)`, galleryNo, user.UserNo, user.Username)
	if e := tx.Error; e != nil {
		db.Rollback()
		return nil, e
	}

	db.Commit()
	return gallery, nil
}

/* Update a Gallery */
func UpdateGallery(cmd *UpdateGalleryCmd, user *util.User) error {

	db := config.GetDB()
	galleryNo := cmd.GalleryNo

	gallery, e := FindGallery(galleryNo)
	if e != nil {
		return e
	}

	// only owner can update the gallery
	if user.UserNo != gallery.UserNo {
		return weberr.NewWebErr("You are not allowed to update this gallery")
	}

	tx := db.Where("gallery_no = ?", galleryNo).Updates(Gallery{
		GalleryNo: cmd.GalleryNo,
		Name:      cmd.Name,
		UpdateBy:  user.Username,
	})

	if e := tx.Error; e != nil {
		log.Warnf("Failed to update gallery, gallery_no: %v, e: %v", galleryNo, tx.Error)
		return weberr.NewWebErr("Failed to update gallery, please try again later")
	}

	return nil
}

/* Find Gallery by gallery_no */
func FindGallery(galleryNo string) (*Gallery, error) {

	db := config.GetDB()
	var gallery Gallery

	tx := db.Raw(`
		SELECT g.* from gallery g 
		WHERE g.gallery_no = ?
		AND g.is_del = 0`, galleryNo).Scan(&gallery)

	if e := tx.Error; e != nil || tx.RowsAffected < 1 {
		if e != nil {
			return nil, tx.Error
		}
		return nil, weberr.NewWebErr("Gallery doesn't exist")
	}
	return &gallery, nil
}

/* Delete a gallery */
func DeleteGallery(cmd *DeleteGalleryCmd, user *util.User) error {

	galleryNo := cmd.GalleryNo
	db := config.GetDB()

	if access, err := HasAccessToGallery(user.UserNo, galleryNo); !access || err != nil {
		if err != nil {
			return err
		}
		return weberr.NewWebErr("You are not allowed to delete this gallery")
	}

	tx := db.Exec(`
		UPDATE gallery g 
		SET g.is_del = 1
		WHERE gallery_no = ? AND g.is_del = 0`, galleryNo)

	if e := tx.Error; e != nil {
		return tx.Error
	}

	return nil
}

// Check if the gallery exists
func GalleryExists(galleryNo string) (bool, error) {

	db := config.GetDB()
	var gallery Gallery

	tx := db.Raw(`
		SELECT g.id from gallery g 
		WHERE g.gallery_no = ?
		AND g.is_del = 0`, galleryNo).Scan(&gallery)

	if e := tx.Error; e != nil || tx.RowsAffected < 1 {
		if e != nil {
			return false, tx.Error
		}
		return false, nil
	}

	return true, nil
}

// Grant user's access to the gallery, only the owner can do so
func GrantGalleryAccessToUser(cmd *PermitGalleryAccessCmd, user *util.User) error {

	gallery, e := FindGallery(cmd.GalleryNo)
	if e != nil {
		return e
	}

	if gallery.UserNo != user.UserNo {
		return weberr.NewWebErr("You are not allowed to grant access to this gallery")
	}

	return CreateGalleryAccess(cmd.UserNo, cmd.GalleryNo, user.Username)
}
