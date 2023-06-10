package controller

import (
	"fmt"

	"github.com/curtisnewbie/fantahsea/client"
	"github.com/curtisnewbie/fantahsea/data"
	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/server"
	"github.com/gin-gonic/gin"
)

/*
	List images of a gallery

	Request Body (JSON): ListGalleryImagesCmd
*/
func ListImagesEndpoint(c *gin.Context, ec common.ExecContext, cmd data.ListGalleryImagesCmd) (any, error) {
	if e := common.Validate(cmd); e != nil {
		return nil, e
	}
	return data.ListGalleryImages(cmd, ec)
}

type TransferGalleryImageReq struct {
	Images []data.CreateGalleryImageCmd
}

/*
	Transfer image from file-server as a gallery image

	Request Body (JSON): TransferGalleryImageReq
*/
func TransferGalleryImageEndpoint(c *gin.Context, ec common.ExecContext, cmd TransferGalleryImageReq) (any, error) {
	user := ec.User
	if e := common.Validate(cmd); e != nil {
		return nil, e
	}

	if cmd.Images == nil || len(cmd.Images) < 1 {
		server.DispatchOk(c)
		return nil, nil
	}

	// validate the keys first
	for _, img := range cmd.Images {
		if isValid, e := client.ValidateFileKey(ec, img.FileKey, user.UserId); e != nil || !isValid {
			if e != nil {
				return nil, e
			}
			return nil, common.NewWebErr(fmt.Sprintf("Only file's owner can make it a gallery image ('%s')", img.Name))
		}
	}

	// start transferring
	go func(ec common.ExecContext, images []data.CreateGalleryImageCmd) {
		for _, cmd := range images {
			fi, er := client.GetFileInfo(ec, cmd.FileKey)
			if er != nil {
				ec.Log.Errorf("Failed to fetch file info while transferring selected images, fi's fileKey: %s, error: %v", cmd.FileKey, er)
				continue
			}

			if fi.Data.FileType == client.FILE { // a file
				if fi.Data.FstoreFileId == "" {
					continue // doesn't have fstore fileId, cannot be transferred
				}

				if data.GuessIsImage(*fi.Data) {
					nc := data.CreateGalleryImageCmd{GalleryNo: cmd.GalleryNo, Name: fi.Data.Name, FileKey: fi.Data.Uuid, FstoreFileId: fi.Data.FstoreFileId}
					if err := data.CreateGalleryImage(ec, nc); err != nil {
						ec.Log.Errorf("Failed to create gallery image, fi's fileKey: %s, error: %v", cmd.FileKey, err)
						continue
					}
				}
			} else { // a directory
				if err := data.TransferImagesInDir(data.TransferGalleryImageInDirReq{
					GalleryNo: cmd.GalleryNo,
					FileKey:   cmd.FileKey,
				}, ec); err != nil {
					ec.Log.Errorf("Failed to transfer images in directory, fi's fileKey: %s, error: %v", cmd.FileKey, err)
					continue
				}
			}
		}
	}(ec, cmd.Images)

	return nil, nil
}
