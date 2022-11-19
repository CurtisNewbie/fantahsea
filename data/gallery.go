package data

import (
	"time"

	"github.com/curtisnewbie/gocommon"
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
	IsDel      gocommon.IS_DEL
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
	Paging    *gocommon.Paging `json:"pagingVo"`
	Galleries *[]VGallery `json:"galleries"`
}

type ListGalleriesCmd struct {
	Paging *gocommon.Paging `json:"pagingVo"`
}

type DeleteGalleryCmd struct {
	GalleryNo string `json:"galleryNo"`
}

type PermitGalleryAccessCmd struct {
	GalleryNo string `json:"galleryNo"`
	UserNo    string `json:"userNo"`
}

type VGalleryBrief struct {
	GalleryNo string `json:"galleryNo"`
	Name      string `json:"name"`
}

type VGallery struct {
	ID         int64     `json:"id"`
	GalleryNo  string    `json:"galleryNo"`
	UserNo     string    `json:"userNo"`
	Name       string    `json:"name"`
	CreateTime gocommon.WTime `json:"createTime"`
	CreateBy   string    `json:"createBy"`
	UpdateTime gocommon.WTime `json:"updateTime"`
	UpdateBy   string    `json:"updateBy"`
	IsOwner    bool      `json:"isOwner"`
}

// List owned gallery briefs
func ListOwnedGalleryBriefs(user *gocommon.User) (*[]VGalleryBrief, error) {
	var briefs []VGalleryBrief
	tx := gocommon.GetMySql().Raw(`select gallery_no, name from gallery 
	where user_no = ? 
	AND is_del = 0`, user.UserNo).Scan(&briefs)

	if e := tx.Error; e != nil {
		return nil, e
	}
	if briefs == nil {
		briefs = []VGalleryBrief{}
	}

	return &briefs, nil
}

/* List Galleries */
func ListGalleries(cmd *ListGalleriesCmd, user *gocommon.User) (*ListGalleriesResp, error) {
	paging := cmd.Paging

	const selectSql string = `
		SELECT g.* from gallery g 
		WHERE (g.user_no = ? 
		OR EXISTS (SELECT * FROM gallery_user_access ga WHERE ga.gallery_no = g.gallery_no AND ga.user_no = ?))
		AND g.is_del = 0 
		LIMIT ?, ?
	`
	db := gocommon.GetMySql()
	var galleries []VGallery

	offset := gocommon.CalcOffset(paging)
	tx := db.Raw(selectSql, user.UserNo, user.UserNo, offset, paging.Limit).Scan(&galleries)

	if e := tx.Error; e != nil {
		return nil, e
	}

	const countSql string = `
		SELECT count(*) from gallery g 
		WHERE (g.user_no = ? 
		OR EXISTS (SELECT * FROM gallery_user_access ga WHERE ga.gallery_no = g.gallery_no AND ga.user_no = ?))
		AND g.is_del = 0
	`
	var total int
	tx = db.Raw(countSql, user.UserNo, user.UserNo).Scan(&total)

	if e := tx.Error; e != nil {
		return nil, e
	}

	if galleries == nil {
		galleries = []VGallery{}
	}

	for i, g := range galleries {
		if g.UserNo == user.UserNo {
			g.IsOwner = true
			galleries[i] = g
		}
	}

	return &ListGalleriesResp{Galleries: &galleries, Paging: gocommon.BuildResPage(paging, total)}, nil
}

// Check if the name is already used by current user
func IsGalleryNameUsed(name string, userNo string) (bool, error) {
	var gallery Gallery
	tx := gocommon.GetMySql().Raw(`
		SELECT g.id from gallery g 
		WHERE g.user_no = ? and g.name = ?
		AND g.is_del = 0`, userNo, name).Scan(&gallery)

	if e := tx.Error; e != nil {
		return false, tx.Error
	}

	return tx.RowsAffected > 0, nil
}

// Create a new Gallery
func CreateGallery(cmd *CreateGalleryCmd, user *gocommon.User) (*Gallery, error) {
	log.Printf("Creating gallery, cmd: %v, user: %v", cmd, user)

	// Guest is not allowed to create gallery
	if gocommon.IsGuest(user) {
		return nil, gocommon.NewWebErr("Guest is not allowed to create gallery")
	}

	if isUsed, err := IsGalleryNameUsed(cmd.Name, user.UserNo); isUsed || err != nil {
		if err != nil {
			return nil, err
		}
		return nil, gocommon.NewWebErr("You already have a gallery with the same name, please change and try again")
	}

	galleryNo := gocommon.GenNoL("GAL", 25)

	db := gocommon.GetMySql().Begin()
	gallery := &Gallery{
		GalleryNo: galleryNo,
		Name:      cmd.Name,
		UserNo:    user.UserNo,
		CreateBy:  user.Username,
		UpdateBy:  user.Username,
		IsDel:     gocommon.IS_DEL_N,
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
func UpdateGallery(cmd *UpdateGalleryCmd, user *gocommon.User) error {

	db := gocommon.GetMySql()
	galleryNo := cmd.GalleryNo

	gallery, e := FindGallery(galleryNo)
	if e != nil {
		return e
	}

	// only owner can update the gallery
	if user.UserNo != gallery.UserNo {
		return gocommon.NewWebErr("You are not allowed to update this gallery")
	}

	tx := db.Where("gallery_no = ?", galleryNo).Updates(Gallery{
		Name:     cmd.Name,
		UpdateBy: user.Username,
	})

	if e := tx.Error; e != nil {
		log.Warnf("Failed to update gallery, gallery_no: %v, e: %v", galleryNo, tx.Error)
		return gocommon.NewWebErr("Failed to update gallery, please try again later")
	}

	return nil
}

/* Find Gallery's creator by gallery_no */
func FindGalleryCreator(galleryNo string) (*string, error) {

	db := gocommon.GetMySql()
	var gallery Gallery

	tx := db.Raw(`
		SELECT g.user_no from gallery g 
		WHERE g.gallery_no = ?
		AND g.is_del = 0`, galleryNo).Scan(&gallery)

	if e := tx.Error; e != nil || tx.RowsAffected < 1 {
		if e != nil {
			return nil, tx.Error
		}
		return nil, gocommon.NewWebErr("Gallery doesn't exist")
	}
	return &gallery.UserNo, nil
}

/* Find Gallery by gallery_no */
func FindGallery(galleryNo string) (*Gallery, error) {

	db := gocommon.GetMySql()
	var gallery Gallery

	tx := db.Raw(`
		SELECT g.* from gallery g 
		WHERE g.gallery_no = ?
		AND g.is_del = 0`, galleryNo).Scan(&gallery)

	if e := tx.Error; e != nil || tx.RowsAffected < 1 {
		if e != nil {
			return nil, tx.Error
		}
		return nil, gocommon.NewWebErr("Gallery doesn't exist")
	}
	return &gallery, nil
}

/* Delete a gallery */
func DeleteGallery(cmd *DeleteGalleryCmd, user *gocommon.User) error {

	galleryNo := cmd.GalleryNo
	db := gocommon.GetMySql()

	if access, err := HasAccessToGallery(user.UserNo, galleryNo); !access || err != nil {
		if err != nil {
			return err
		}
		return gocommon.NewWebErr("You are not allowed to delete this gallery")
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

	db := gocommon.GetMySql()
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
func GrantGalleryAccessToUser(cmd *PermitGalleryAccessCmd, user *gocommon.User) error {

	gallery, e := FindGallery(cmd.GalleryNo)
	if e != nil {
		return e
	}

	if gallery.UserNo != user.UserNo {
		return gocommon.NewWebErr("You are not allowed to grant access to this gallery")
	}

	return CreateGalleryAccess(cmd.UserNo, cmd.GalleryNo, user.Username)
}
