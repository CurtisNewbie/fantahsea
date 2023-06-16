package data

import (
	"fmt"
	"strconv"
	"time"

	"github.com/curtisnewbie/fantahsea/client"
	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/mysql"
)

// GalleryImage.status (doesn't really matter anymore)
type ImgStatus string

const (
	NORMAL  ImgStatus = "NORMAL"
	DELETED ImgStatus = "DELETED"

	// 40mb is the maximum size for an image
	IMAGE_SIZE_THRESHOLD int64 = 40 * 1048576
)

// ------------------------------- entity start

type TransferGalleryImageInDirReq struct {
	// gallery no
	GalleryNo string `json:"galleryNo" validation:"notEmpty"`

	// file key of the directory
	FileKey string `json:"fileKey" validation:"notEmpty"`
}

// Image that belongs to a Gallery
type GalleryImage struct {
	ID         int64
	GalleryNo  string
	ImageNo    string
	Name       string
	FileKey    string
	Status     ImgStatus
	CreateTime time.Time
	CreateBy   string
	UpdateTime time.Time
	UpdateBy   string
	IsDel      common.IS_DEL
}

func (GalleryImage) TableName() string {
	return "gallery_image"
}

// ------------------------------- entity end

type ThumbnailInfo struct {
	Name string
	Path string
}

type ListGalleryImagesCmd struct {
	GalleryNo     string `json:"galleryNo" validation:"notEmpty"`
	common.Paging `json:"pagingVo"`
}

type ListGalleryImagesResp struct {
	Images []ImageInfo   `json:"images"`
	Paging common.Paging `json:"pagingVo"`
}

type ImageInfo struct {
	ImageFileId     string `json:"-"`
	ThumbnailFileId string `json:"-"`
	fileKey         string `json:"-"`
	ThumbnailToken  string `json:"thumbnailToken"`
	FileTempToken   string `json:"fileTempToken"`
}

type CreateGalleryImageCmd struct {
	GalleryNo    string `json:"galleryNo"`
	Name         string `json:"name"`
	FileKey      string `json:"fileKey"`
	FstoreFileId string
}

// Create a gallery image record
func CreateGalleryImage(ec common.ExecContext, cmd CreateGalleryImageCmd) error {
	user := ec.User
	creator, err := FindGalleryCreator(cmd.GalleryNo)
	if err != nil {
		return err
	}

	if *creator != user.UserNo {
		return common.NewWebErr("You are not allowed to upload image to this gallery")
	}

	if isCreated, e := isImgCreatedAlready(cmd.GalleryNo, cmd.FileKey); isCreated || e != nil {
		if e != nil {
			return e
		}
		ec.Log.Infof("Image '%s' added already", cmd.Name)
		return nil
	}

	imageNo := common.GenNoL("IMG", 25)
	const sql string = `
			insert into gallery_image (gallery_no, image_no, name, file_key, create_by)
			values (?, ?, ?, ?, ?)
		`
	return mysql.GetConn().Exec(sql, cmd.GalleryNo, imageNo, cmd.Name, cmd.FileKey, user.Username).Error
}

// List gallery images
func ListGalleryImages(cmd ListGalleryImagesCmd, ec common.ExecContext) (*ListGalleryImagesResp, error) {
	user := ec.User
	ec.Log.Infof("ListGalleryImages, cmd: %+v", cmd)

	if hasAccess, err := HasAccessToGallery(user.UserNo, cmd.GalleryNo); err != nil || !hasAccess {
		if err != nil {
			return nil, err
		}
		return nil, common.NewWebErr("You are not allowed to access this gallery")
	}

	const selectSql string = `
		select image_no, file_key from gallery_image
		where gallery_no = ?
		order by id desc
		limit ?, ?
	`
	var galleryImages []GalleryImage
	tx := mysql.GetMySql().Raw(selectSql, cmd.GalleryNo, cmd.Paging.GetOffset(), cmd.Paging.GetLimit()).Scan(&galleryImages)
	if tx.Error != nil {
		return nil, tx.Error
	}

	if galleryImages == nil {
		galleryImages = []GalleryImage{}
	}

	// count total asynchronoulsy (normally, when the SELECT is successful, the COUNT doesn't really fail)
	countFuture := common.RunAsync(func() (int, error) {
		var total int
		tx = mysql.GetMySql().
			Raw(`select count(*) from gallery_image where gallery_no = ?`, cmd.GalleryNo).
			Scan(&total)
		return total, tx.Error
	})

	// generate temp tokens for the actual files and the thumbnail, these are served by mini-fstore
	images := []ImageInfo{}
	if len(galleryImages) > 0 {
		genTknReqs := []client.BatchGenFileKeyItem{}
		for _, img := range galleryImages {
			r, e := client.GetFileInfo(ec, img.FileKey)
			if e != nil {
				return nil, e
			}
			fstoreFileId := r.Data.FstoreFileId
			genTknReqs = append(genTknReqs, client.BatchGenFileKeyItem{FileId: fstoreFileId, Filename: r.Data.Name})

			thumbnailFileId := r.Data.Thumbnail
			if thumbnailFileId == "" {
				thumbnailFileId = fstoreFileId
			} else {
				genTknReqs = append(genTknReqs, client.BatchGenFileKeyItem{FileId: thumbnailFileId, Filename: r.Data.Name})
			}

			images = append(images, ImageInfo{fileKey: img.FileKey, ImageFileId: fstoreFileId, ThumbnailFileId: thumbnailFileId})
		}

		// requests temp tokens in batch
		tokens, err := client.BatchGetFstoreTmpToken(ec, client.BatchGenFileKeyReq{Items: genTknReqs})
		if err != nil {
			return nil, fmt.Errorf("failed to generate fstore temp tokens in batch, %v", err)
		}

		idTknMap := map[string]string{}
		for _, t := range tokens {
			idTknMap[t.FileId] = t.TempKey
		}
		for i, im := range images {
			im.ThumbnailToken = idTknMap[im.ThumbnailFileId]
			im.FileTempToken = idTknMap[im.ImageFileId]
			images[i] = im
		}
	}

	total, errCnt := countFuture.Get()
	if errCnt != nil {
		return nil, errCnt
	}

	return &ListGalleryImagesResp{Images: images, Paging: common.RespPage(cmd.Paging, total)}, nil
}

// Transfer images in dir
func TransferImagesInDir(cmd TransferGalleryImageInDirReq, ec common.ExecContext) error {
	user := ec.User
	resp, e := client.GetFileInfo(ec, cmd.FileKey)
	if e != nil {
		return e
	}

	fi := resp.Data

	// only the owner of the directory can do this, by default directory is only visible to the uploader
	if strconv.Itoa(fi.UploaderId) != user.UserId {
		return common.NewWebErr("Not permitted operation")
	}

	if fi.FileType != client.DIR {
		return common.NewWebErr("This is not a directory")
	}

	if fi.IsDeleted {
		return common.NewWebErr("Directory is already deleted")
	}
	dirFileKey := cmd.FileKey
	galleryNo := cmd.GalleryNo
	start := time.Now()

	page := 1
	for {
		resp, err := client.ListFilesInDir(ec, dirFileKey, 100, page)
		if err != nil {
			ec.Log.Errorf("Failed to list files in dir, dir's fileKey: %s, error: %v", dirFileKey, err)
			break
		}
		if resp.Data == nil || len(resp.Data) < 1 {
			break
		}

		// starts fetching file one by one
		for i := 0; i < len(resp.Data); i++ {
			fk := resp.Data[i]
			fi, er := client.GetFileInfo(ec, fk)
			if er != nil {
				ec.Log.Errorf("Failed to fetch file info while looping files in dir, fi's fileKey: %s, error: %v", fk, er)
				continue
			}

			if GuessIsImage(ec, *fi.Data) {
				cmd := CreateGalleryImageCmd{GalleryNo: galleryNo, Name: fi.Data.Name, FileKey: fi.Data.Uuid, FstoreFileId: fi.Data.FstoreFileId}
				if err := CreateGalleryImage(ec, cmd); err != nil {
					ec.Log.Errorf("Failed to create gallery image, fi's fileKey: %s, error: %v", fk, err)
				}
			}
		}

		page += 1
	}

	ec.Log.Infof("Finished TransferImagesInDir, dir's fileKey: %s, took: %s", dirFileKey, time.Since(start))
	return nil
}

/*
	-----------------------------------------------------------

	Helper methods

	-----------------------------------------------------------
*/

// Guess whether a file is an image
func GuessIsImage(c common.ExecContext, f client.FileInfoResp) bool {
	if f.SizeInBytes > IMAGE_SIZE_THRESHOLD {
		return false
	}
	if f.FileType != client.FILE {
		return false
	}
	if f.Thumbnail == "" {
		c.Log.Infof("File doesn't have thumbnail, fileKey: %v", f.Uuid)
		return false
	}

	return true
}

// check whether the gallery image is created already
//
// return isImgCreated, error
func isImgCreatedAlready(galleryNo string, fileKey string) (bool, error) {
	db := mysql.GetMySql()

	var id int
	tx := db.Raw(`
		SELECT id FROM gallery_image
		WHERE gallery_no = ?
		AND file_key = ?
		AND is_del = 0
		`, galleryNo, fileKey).Scan(&id)

	if e := tx.Error; e != nil || tx.RowsAffected < 1 {
		return false, tx.Error
	}

	return true, nil
}
