package fantahsea

import (
	"fmt"
	"time"

	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/miso/miso"
)

// GalleryImage.status (doesn't really matter anymore)
type ImgStatus string

const (
	NORMAL  ImgStatus = "NORMAL"
	DELETED ImgStatus = "DELETED"

	// 40mb is the maximum size for an image
	IMAGE_SIZE_THRESHOLD int64 = 40 * 1048576
)

type TransferGalleryImageReq struct {
	Images []CreateGalleryImageCmd
}

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

type ThumbnailInfo struct {
	Name string
	Path string
}

type CreateGalleryImgEvent struct {
	Username     string `json:"username"`
	UserNo       string `json:"userNo"`
	DirFileKey   string `json:"dirFileKey"`
	DirName      string `json:"dirName"`
	ImageName    string `json:"imageName"`
	ImageFileKey string `json:"imageFileKey"`
}

type NotifyFileDeletedEvent struct {
	FileKey string `json:"fileKey"`
}

type ListGalleryImagesCmd struct {
	GalleryNo   string `json:"galleryNo" validation:"notEmpty"`
	miso.Paging `json:"pagingVo"`
}

type ListGalleryImagesResp struct {
	Images []ImageInfo `json:"images"`
	Paging miso.Paging `json:"pagingVo"`
}

type ImageInfo struct {
	ThumbnailToken  string `json:"thumbnailToken"`
	FileTempToken   string `json:"fileTempToken"`
	ImageFileId     string `json:"-"`
	ThumbnailFileId string `json:"-"`
}

type CreateGalleryImageCmd struct {
	GalleryNo string `json:"galleryNo"`
	Name      string `json:"name"`
	FileKey   string `json:"fileKey"`
}

func DeleteGalleryImage(rail miso.Rail, fileKey string) error {
	return miso.GetMySQL().Exec("delete from gallery_image where file_key = ?", fileKey).Error
}

// Create a gallery image record
func CreateGalleryImage(rail miso.Rail, cmd CreateGalleryImageCmd, userNo string, username string) error {
	creator, err := FindGalleryCreator(rail, cmd.GalleryNo)
	if err != nil {
		return err
	}

	if *creator != userNo {
		return miso.NewErr("You are not allowed to upload image to this gallery")
	}

	if isCreated, e := isImgCreatedAlready(rail, cmd.GalleryNo, cmd.FileKey); isCreated || e != nil {
		if e != nil {
			return e
		}
		rail.Infof("Image '%s' added already", cmd.Name)
		return nil
	}

	imageNo := miso.GenNoL("IMG", 25)
	const sql string = `
			insert into gallery_image (gallery_no, image_no, name, file_key, create_by)
			values (?, ?, ?, ?, ?)
		`
	return miso.GetMySQL().Exec(sql, cmd.GalleryNo, imageNo, cmd.Name, cmd.FileKey, username).Error
}

// List gallery images
func ListGalleryImages(rail miso.Rail, cmd ListGalleryImagesCmd, user common.User) (*ListGalleryImagesResp, error) {
	rail.Infof("ListGalleryImages, cmd: %+v", cmd)

	if hasAccess, err := HasAccessToGallery(user.UserNo, cmd.GalleryNo); err != nil || !hasAccess {
		if err != nil {
			return nil, err
		}
		return nil, miso.NewErr("You are not allowed to access this gallery")
	}

	const selectSql string = `
		select image_no, file_key from gallery_image
		where gallery_no = ?
		order by id desc
		limit ?, ?
	`
	var galleryImages []GalleryImage
	tx := miso.GetMySQL().Raw(selectSql, cmd.GalleryNo, cmd.Paging.GetOffset(), cmd.Paging.GetLimit()).Scan(&galleryImages)
	if tx.Error != nil {
		return nil, tx.Error
	}

	if galleryImages == nil {
		galleryImages = []GalleryImage{}
	}

	// count total asynchronoulsy (normally, when the SELECT is successful, the COUNT doesn't really fail)
	countFuture := miso.RunAsync(func() (int, error) {
		var total int
		tx = miso.GetMySQL().
			Raw(`select count(*) from gallery_image where gallery_no = ?`, cmd.GalleryNo).
			Scan(&total)
		return total, tx.Error
	})

	// generate temp tokens for the actual files and the thumbnail, these are served by mini-fstore
	images := []ImageInfo{}
	if len(galleryImages) > 0 {

		genTknFutures := []miso.Future[BatchGenFileKeyResp]{}

		for _, img := range galleryImages {
			r, e := GetFileInfo(rail, img.FileKey)
			if e != nil {
				rail.Errorf("GetFileInfo failed, fileKey: %v, %v", img.FileKey, e)
				continue
			}
			fstoreFileId := r.Data.FstoreFileId

			genTknFutures = append(genTknFutures, miso.RunAsync[BatchGenFileKeyResp](func() (BatchGenFileKeyResp, error) {
				tkn, err := GetFstoreTmpToken(rail.NextSpan(), fstoreFileId, r.Data.Name)
				if err != nil {
					return BatchGenFileKeyResp{FileId: fstoreFileId}, err
				}
				return BatchGenFileKeyResp{
					FileId:  fstoreFileId,
					TempKey: tkn,
				}, nil
			}))

			thumbnailFileId := r.Data.Thumbnail
			if thumbnailFileId == "" {
				thumbnailFileId = fstoreFileId
			} else {
				genTknFutures = append(genTknFutures, miso.RunAsync[BatchGenFileKeyResp](func() (BatchGenFileKeyResp, error) {
					tkn, err := GetFstoreTmpToken(rail.NextSpan(), thumbnailFileId, r.Data.Name)
					if err != nil {
						return BatchGenFileKeyResp{FileId: thumbnailFileId}, err
					}
					return BatchGenFileKeyResp{
						FileId:  thumbnailFileId,
						TempKey: tkn,
					}, nil
				}))
			}

			images = append(images, ImageInfo{ImageFileId: fstoreFileId, ThumbnailFileId: thumbnailFileId})
		}

		tokens := []BatchGenFileKeyResp{}

		for i := range genTknFutures {
			res, err := genTknFutures[i].Get()
			if err != nil {
				rail.Errorf("Failed to get mini-fstore temp token for fstore_file_id: %v, %v", res.FileId, err)
				continue
			}
			tokens = append(tokens, res)
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

	return &ListGalleryImagesResp{Images: images, Paging: miso.RespPage(cmd.Paging, total)}, nil
}

func BatchTransferAsync(rail miso.Rail, cmd TransferGalleryImageReq, user common.User) (any, error) {
	if cmd.Images == nil || len(cmd.Images) < 1 {
		return nil, nil
	}

	// validate the keys first
	for _, img := range cmd.Images {
		if isValid, e := ValidateFileKey(rail, img.FileKey, user.UserId); e != nil || !isValid {
			if e != nil {
				return nil, e
			}
			return nil, miso.NewErr(fmt.Sprintf("Only file's owner can make it a gallery image ('%s')", img.Name))
		}
	}

	// start transferring
	go func(rail miso.Rail, images []CreateGalleryImageCmd) {
		for _, cmd := range images {
			fi, er := GetFileInfo(rail, cmd.FileKey)
			if er != nil {
				rail.Errorf("Failed to fetch file info while transferring selected images, fi's fileKey: %s, error: %v", cmd.FileKey, er)
				continue
			}

			if fi.Data.FileType == FILE { // a file
				if fi.Data.FstoreFileId == "" {
					continue // doesn't have fstore fileId, cannot be transferred
				}

				if GuessIsImage(rail, *fi.Data) {
					nc := CreateGalleryImageCmd{GalleryNo: cmd.GalleryNo, Name: fi.Data.Name, FileKey: fi.Data.Uuid}
					if err := CreateGalleryImage(rail, nc, user.UserNo, user.Username); err != nil {
						rail.Errorf("Failed to create gallery image, fi's fileKey: %s, error: %v", cmd.FileKey, err)
						continue
					}
				}
			} else { // a directory
				treq := TransferGalleryImageInDirReq{
					GalleryNo: cmd.GalleryNo,
					FileKey:   cmd.FileKey,
				}
				if err := TransferImagesInDir(rail, treq, user); err != nil {
					rail.Errorf("Failed to transfer images in directory, fi's fileKey: %s, error: %v", cmd.FileKey, err)
					continue
				}
			}
		}
	}(rail, cmd.Images)

	return nil, nil
}

// Transfer images in dir
func TransferImagesInDir(rail miso.Rail, cmd TransferGalleryImageInDirReq, user common.User) error {
	resp, e := GetFileInfo(rail, cmd.FileKey)
	if e != nil {
		return e
	}

	fi := resp.Data

	// only the owner of the directory can do this, by default directory is only visible to the uploader
	if fi.UploaderId != user.UserId {
		return miso.NewErr("Not permitted operation")
	}

	if fi.FileType != DIR {
		return miso.NewErr("This is not a directory")
	}

	if fi.IsDeleted {
		return miso.NewErr("Directory is already deleted")
	}
	dirFileKey := cmd.FileKey
	galleryNo := cmd.GalleryNo
	start := time.Now()

	page := 1
	for {
		resp, err := ListFilesInDir(rail, dirFileKey, 100, page)
		if err != nil {
			rail.Errorf("Failed to list files in dir, dir's fileKey: %s, error: %v", dirFileKey, err)
			break
		}
		if resp.Data == nil || len(resp.Data) < 1 {
			break
		}

		// starts fetching file one by one
		for i := 0; i < len(resp.Data); i++ {
			fk := resp.Data[i]
			fi, er := GetFileInfo(rail, fk)
			if er != nil {
				rail.Errorf("Failed to fetch file info while looping files in dir, fi's fileKey: %s, error: %v", fk, er)
				continue
			}

			if GuessIsImage(rail, *fi.Data) {
				cmd := CreateGalleryImageCmd{GalleryNo: galleryNo, Name: fi.Data.Name, FileKey: fi.Data.Uuid}
				if err := CreateGalleryImage(rail, cmd, user.UserNo, user.Username); err != nil {
					rail.Errorf("Failed to create gallery image, fi's fileKey: %s, error: %v", fk, err)
				}
			}
		}

		page += 1
	}

	rail.Infof("Finished TransferImagesInDir, dir's fileKey: %s, took: %s", dirFileKey, time.Since(start))
	return nil
}

// Guess whether a file is an image
func GuessIsImage(rail miso.Rail, f FileInfoResp) bool {
	if f.SizeInBytes > IMAGE_SIZE_THRESHOLD {
		return false
	}
	if f.FileType != FILE {
		return false
	}
	if f.Thumbnail == "" {
		rail.Infof("File doesn't have thumbnail, fileKey: %v", f.Uuid)
		return false
	}

	return true
}

// check whether the gallery image is created already
//
// return isImgCreated, error
func isImgCreatedAlready(rail miso.Rail, galleryNo string, fileKey string) (bool, error) {
	db := miso.GetMySQL()

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
