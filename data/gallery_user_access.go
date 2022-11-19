package data

import (
	"time"

	"github.com/curtisnewbie/gocommon"
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
	IsDel      gocommon.IS_DEL
}

type UpdateGUAIsDelCmd struct {
	GalleryNo string
	UserNo    string
	IsDelFrom gocommon.IS_DEL
	IsDelTo   gocommon.IS_DEL
	UpdateBy  string
}

func (GalleryUserAccess) TableName() string {
	return "gallery_user_access"
}

// ------------------------------- entity end

/* Check if user has access to the gallery */
func HasAccessToGallery(userNo string, galleryNo string) (bool, error) {

	// check if the user has access to the gallery
	userAccess, err := findGalleryAccess(userNo, galleryNo)
	if err != nil {
		return false, err
	}

	if userAccess == nil || gocommon.IsDeleted(userAccess.IsDel) {
		return false, nil
	}

	return true, nil
}

// Assign user access to the gallery
func CreateGalleryAccess(userNo string, galleryNo string, operator string) error {

	// check if the user has access to the gallery
	userAccess, err := findGalleryAccess(userNo, galleryNo)
	if err != nil {
		return err
	}

	if userAccess != nil && !gocommon.IsDeleted(userAccess.IsDel) {
		return nil
	}

	var e error
	if userAccess == nil {
		e = createUserAccess(userNo, galleryNo, operator)
	} else {
		e = updateUserAccessIsDelFlag(&UpdateGUAIsDelCmd{
			UserNo:    userNo,
			GalleryNo: galleryNo,
			IsDelFrom: gocommon.IS_DEL_N,
			IsDelTo:   gocommon.IS_DEL_Y,
			UpdateBy:  operator,
		})
	}

	return e
}

/*
	-----------------------------------------------------------

	Helper methods

	-----------------------------------------------------------
*/

/* find GalleryUserAccess, is_del flag is ignored */
func findGalleryAccess(userNo string, galleryNo string) (*GalleryUserAccess, error) {

	db := gocommon.GetMySql()

	// check if the user has access to the gallery
	var userAccess *GalleryUserAccess = &GalleryUserAccess{}

	tx := db.Raw(`
		SELECT * FROM gallery_user_access 
		WHERE gallery_no = ?
		AND user_no = ?`, galleryNo, userNo).Scan(&userAccess)

	if e := tx.Error; e != nil || tx.RowsAffected < 1 {
		if e != nil {
			return nil, e
		}
		return nil, nil
	}

	return userAccess, nil
}

// Insert a new gallery_user_access record
func createUserAccess(userNo string, galleryNo string, createdBy string) error {

	db := gocommon.GetMySql()

	tx := db.Exec(`INSERT INTO gallery_user_access (gallery_no, user_no, create_by) VALUES (?, ?, ?)`, galleryNo, userNo, createdBy)

	if e := tx.Error; e != nil {
		return e
	}

	return nil
}

// Update is_del of the record
func updateUserAccessIsDelFlag(cmd *UpdateGUAIsDelCmd) error {

	tx := gocommon.GetMySql().Exec(`
	UPDATE gallery_user_access SET is_del = ?, update_by = ?
	WHERE gallery_no = ? AND user_no = ? AND is_del = ?`, cmd.IsDelTo, cmd.UpdateBy, cmd.GalleryNo, cmd.UserNo, cmd.IsDelFrom)

	if e := tx.Error; e != nil {
		return e
	}

	return nil
}
