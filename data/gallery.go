package data

import (
	"time"

	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/mysql"
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
	IsDel      common.IS_DEL
}

func (Gallery) TableName() string {
	return "gallery"

}

// ------------------------------- entity end

type CreateGalleryCmd struct {
	Name string `json:"name" validation:"notEmpty"`
}

type UpdateGalleryCmd struct {
	GalleryNo string `json:"galleryNo" validation:"notEmpty"`
	Name      string `json:"name" validation:"notEmpty"`
}

type ListGalleriesResp struct {
	Paging    common.Paging `json:"pagingVo"`
	Galleries []VGallery    `json:"galleries"`
}

type ListGalleriesCmd struct {
	Paging common.Paging `json:"pagingVo"`
}

type DeleteGalleryCmd struct {
	GalleryNo string `json:"galleryNo" validation:"notEmpty"`
}

type PermitGalleryAccessCmd struct {
	GalleryNo string `json:"galleryNo" validation:"notEmpty"`
	UserNo    string `json:"userNo" validation:"notEmpty"`
}

type VGalleryBrief struct {
	GalleryNo string `json:"galleryNo"`
	Name      string `json:"name"`
}

type VGallery struct {
	ID         int64        `json:"id"`
	GalleryNo  string       `json:"galleryNo"`
	UserNo     string       `json:"userNo"`
	Name       string       `json:"name"`
	CreateTime common.WTime `json:"createTime"`
	CreateBy   string       `json:"createBy"`
	UpdateTime common.WTime `json:"updateTime"`
	UpdateBy   string       `json:"updateBy"`
	IsOwner    bool         `json:"isOwner"`
}

// List owned gallery briefs
func ListOwnedGalleryBriefs(ec common.ExecContext) (*[]VGalleryBrief, error) {
	user := ec.User
	var briefs []VGalleryBrief
	tx := mysql.
		GetMySql().
		Raw(`select gallery_no, name from gallery where user_no = ? AND is_del = 0`, user.UserNo).
		Scan(&briefs)

	if e := tx.Error; e != nil {
		return nil, e
	}
	if briefs == nil {
		briefs = []VGalleryBrief{}
	}

	return &briefs, nil
}

/* List Galleries */
func ListGalleries(cmd ListGalleriesCmd, ec common.ExecContext) (ListGalleriesResp, error) {
	paging := cmd.Paging

	const selectSql string = `
		SELECT g.* from gallery g
		WHERE (g.user_no = ?
		OR EXISTS (SELECT * FROM gallery_user_access ga WHERE ga.gallery_no = g.gallery_no AND ga.user_no = ?))
		AND g.is_del = 0
		ORDER BY id DESC
		LIMIT ?, ?
	`
	db := mysql.GetMySql()
	var galleries []VGallery

	offset := paging.GetOffset()
	tx := db.Raw(selectSql, ec.User.UserNo, ec.User.UserNo, offset, paging.Limit).Scan(&galleries)

	if e := tx.Error; e != nil {
		return ListGalleriesResp{}, e
	}

	const countSql string = `
		SELECT count(*) from gallery g
		WHERE (g.user_no = ?
		OR EXISTS (SELECT * FROM gallery_user_access ga WHERE ga.gallery_no = g.gallery_no AND ga.user_no = ?))
		AND g.is_del = 0
	`
	var total int
	tx = db.Raw(countSql, ec.User.UserNo, ec.User.UserNo).Scan(&total)

	if e := tx.Error; e != nil {
		return ListGalleriesResp{}, e
	}

	if galleries == nil {
		galleries = []VGallery{}
	}

	for i, g := range galleries {
		if g.UserNo == ec.User.UserNo {
			g.IsOwner = true
			galleries[i] = g
		}
	}

	return ListGalleriesResp{Galleries: galleries, Paging: paging.ToRespPage(total)}, nil
}

// Check if the name is already used by current user
func IsGalleryNameUsed(name string, userNo string) (bool, error) {
	var gallery Gallery
	tx := mysql.
		GetMySql().
		Raw(`SELECT g.id from gallery g WHERE g.user_no = ? and g.name = ? AND g.is_del = 0`, userNo, name).
		Scan(&gallery)

	if e := tx.Error; e != nil {
		return false, tx.Error
	}

	return tx.RowsAffected > 0, nil
}

// Create a new Gallery
func CreateGallery(cmd CreateGalleryCmd, ec common.ExecContext) (*Gallery, error) {
	user := ec.User
	ec.Log.Infof("Creating gallery, cmd: %v, user: %v", cmd, user)

	// Guest is not allowed to create gallery
	if common.IsGuest(user) {
		return nil, common.NewWebErr("Guest is not allowed to create gallery")
	}

	if isUsed, err := IsGalleryNameUsed(cmd.Name, user.UserNo); isUsed || err != nil {
		if err != nil {
			return nil, err
		}
		return nil, common.NewWebErr("You already have a gallery with the same name, please change and try again")
	}

	galleryNo := common.GenNoL("GAL", 25)

	db := mysql.GetMySql().Begin()
	gallery := &Gallery{
		GalleryNo: galleryNo,
		Name:      cmd.Name,
		UserNo:    user.UserNo,
		CreateBy:  user.Username,
		UpdateBy:  user.Username,
		IsDel:     common.IS_DEL_N,
	}

	result := db.Omit("CreateTime", "UpdateTime").Create(gallery)
	if result.Error != nil {
		db.Rollback()
		return nil, result.Error
	}

	tx := db.Exec(`INSERT INTO gallery_user_access (gallery_no, user_no, create_by) VALUES (?, ?, ?)`, galleryNo, user.UserNo, user.Username)
	if e := tx.Error; e != nil {
		ec.Log.Errorf("Failed to create gallery user access, galleryNo: %s, userNo: %s, username: %s", galleryNo, user.UserNo, user.Username)
		db.Rollback()
		return nil, e
	}

	db.Commit()
	return gallery, nil
}

/* Update a Gallery */
func UpdateGallery(cmd UpdateGalleryCmd, ec common.ExecContext) error {
	user := ec.User
	db := mysql.GetMySql()
	galleryNo := cmd.GalleryNo

	gallery, e := FindGallery(galleryNo)
	if e != nil {
		return e
	}

	// only owner can update the gallery
	if user.UserNo != gallery.UserNo {
		return common.NewWebErr("You are not allowed to update this gallery")
	}

	tx := db.Where("gallery_no = ?", galleryNo).Updates(Gallery{
		Name:     cmd.Name,
		UpdateBy: user.Username,
	})

	if e := tx.Error; e != nil {
		ec.Log.Warnf("Failed to update gallery, gallery_no: %v, e: %v", galleryNo, tx.Error)
		return common.NewWebErr("Failed to update gallery, please try again later")
	}

	return nil
}

/* Find Gallery's creator by gallery_no */
func FindGalleryCreator(galleryNo string) (*string, error) {

	db := mysql.GetMySql()
	var gallery Gallery

	tx := db.Raw(`
		SELECT g.user_no from gallery g
		WHERE g.gallery_no = ?
		AND g.is_del = 0`, galleryNo).Scan(&gallery)

	if e := tx.Error; e != nil || tx.RowsAffected < 1 {
		if e != nil {
			return nil, tx.Error
		}
		return nil, common.NewWebErr("Gallery doesn't exist")
	}
	return &gallery.UserNo, nil
}

/* Find Gallery by gallery_no */
func FindGallery(galleryNo string) (*Gallery, error) {

	db := mysql.GetMySql()
	var gallery Gallery

	tx := db.Raw(`
		SELECT g.* from gallery g
		WHERE g.gallery_no = ?
		AND g.is_del = 0`, galleryNo).Scan(&gallery)

	if e := tx.Error; e != nil || tx.RowsAffected < 1 {
		if e != nil {
			return nil, tx.Error
		}
		return nil, common.NewWebErr("Gallery doesn't exist")
	}
	return &gallery, nil
}

/* Delete a gallery */
func DeleteGallery(cmd DeleteGalleryCmd, ec common.ExecContext) error {
	user := ec.User
	galleryNo := cmd.GalleryNo
	db := mysql.GetMySql()

	if access, err := HasAccessToGallery(user.UserNo, galleryNo); !access || err != nil {
		if err != nil {
			return err
		}
		return common.NewWebErr("You are not allowed to delete this gallery")
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

	db := mysql.GetMySql()
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
func GrantGalleryAccessToUser(cmd PermitGalleryAccessCmd, ec common.ExecContext) error {
	user := ec.User
	gallery, e := FindGallery(cmd.GalleryNo)
	if e != nil {
		return e
	}

	if gallery.UserNo != user.UserNo {
		return common.NewWebErr("You are not allowed to grant access to this gallery")
	}

	return CreateGalleryAccess(cmd.UserNo, cmd.GalleryNo, user.Username)
}
