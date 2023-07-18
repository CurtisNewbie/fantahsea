package data

import (
	"time"

	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/mysql"
	"github.com/curtisnewbie/gocommon/redis"
	"gorm.io/gorm"
)

// ------------------------------- entity start

// Gallery
type Gallery struct {
	ID         int64
	GalleryNo  string
	UserNo     string
	Name       string
	DirFileKey string
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

type CreateGalleryForDirCmd struct {
	DirName    string
	DirFileKey string
	Username   string
	UserNo     string
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

const (
	AddDirGalleryImageEventBus = "fantahsea.dir.gallery.image.add"
	NotifyFileDeletedEventBus  = "fantahsea.notify.file.deleted"
)

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

func GalleryNoOfDir(dirFileKey string) (string, error) {
	var gallery Gallery
	tx := mysql.
		GetMySql().
		Raw(`SELECT g.gallery_no from gallery g WHERE g.dir_file_key = ? and g.is_del = 0 limit 1`, dirFileKey).
		Scan(&gallery)

	if e := tx.Error; e != nil {
		return "", tx.Error
	}

	return gallery.GalleryNo, nil
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

// Create a new Gallery for dir
func CreateGalleryForDir(ec common.ExecContext, cmd CreateGalleryForDirCmd) (string, error) {

	return redis.RLockRun[string](ec, "fantahsea:gallery:create:"+cmd.UserNo,
		func() (string, error) {
			galleryNo, err := GalleryNoOfDir(cmd.DirFileKey)
			if err != nil {
				return "", err
			}

			if galleryNo == "" {
				galleryNo = common.GenNoL("GAL", 25)
				ec.Log.Infof("Creating gallery (%s) for directory %s (%s)", galleryNo, cmd.DirName, cmd.DirFileKey)

				err := mysql.GetConn().Transaction(func(tx *gorm.DB) error {
					gallery := &Gallery{
						GalleryNo:  galleryNo,
						Name:       cmd.DirName,
						DirFileKey: cmd.DirFileKey,
						UserNo:     cmd.UserNo,
						CreateBy:   cmd.Username,
						UpdateBy:   cmd.Username,
						IsDel:      common.IS_DEL_N,
					}

					result := tx.Omit("CreateTime", "UpdateTime").Create(gallery)
					if result.Error != nil {
						return result.Error
					}

					t := tx.Exec(`INSERT INTO gallery_user_access (gallery_no, user_no, create_by) VALUES (?, ?, ?)`, galleryNo, cmd.UserNo, cmd.Username)
					if e := t.Error; e != nil {
						ec.Log.Errorf("Failed to create gallery user access, galleryNo: %s, userNo: %s, username: %s", galleryNo, cmd.UserNo, cmd.Username)
						return e
					}

					return nil
				})
				if err != nil {
					return galleryNo, err
				}
			}
			return galleryNo, nil
		})
}

// Create a new Gallery
func CreateGallery(ec common.ExecContext, cmd CreateGalleryCmd) (*Gallery, error) {
	user := ec.User
	ec.Log.Infof("Creating gallery, cmd: %v, user: %v", cmd, user)

	gal, er := redis.RLockRun(ec, "fantahsea:gallery:create:"+ec.User.UserNo, func() (*Gallery, error) {

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

	})

	if er != nil {
		return nil, er
	}

	return gal, nil
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

func OnCreateGalleryImgEvent(c common.ExecContext, evt CreateGalleryImgEvent) error {

	// it's meant to be used for adding image to the gallery that belongs to the directory
	if evt.DirFileKey == "" {
		return nil
	}

	// create gallery for the directory if necessary
	galleryNo, err := CreateGalleryForDir(c, CreateGalleryForDirCmd{
		Username:   evt.Username,
		UserNo:     evt.UserNo,
		DirName:    evt.DirName,
		DirFileKey: evt.DirFileKey,
	})

	if err != nil {
		return err
	}

	// add image to the gallery
	return CreateGalleryImage(c,
		CreateGalleryImageCmd{
			GalleryNo: galleryNo,
			Name:      evt.ImageName,
			FileKey:   evt.ImageFileKey,
		},
		evt.UserNo,
		evt.Username)
}

func OnNotifyFileDeletedEvent(c common.ExecContext, evt NotifyFileDeletedEvent) error {
	c.Log.Infof("Received NotifyFileDeletedEvent, fileKey: %v", evt.FileKey)
	return DeleteGalleryImage(c, evt.FileKey)
}
