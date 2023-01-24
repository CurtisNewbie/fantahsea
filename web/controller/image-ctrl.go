package controller

import (
	"fmt"
	"net/http"

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
func ListImagesEndpoint(c *gin.Context, ec common.ExecContext) (any, error) {
	var cmd data.ListGalleryImagesCmd
	server.MustBindJson(c, &cmd)
	if e := common.Validate(cmd); e != nil {
		return nil, e
	}
	return data.ListGalleryImages(cmd, ec)
}

/*
	Download image thumbnail
*/
func DownloadImageThumbnailEndpoint(c *gin.Context, ec common.ExecContext) {
	token := c.Query("token")
	ec.Log.Printf("Download Image thumbnail, token: %s", token)
	dimg, e := data.ResolveImageThumbnail(ec, token)
	if e != nil {
		ec.Log.Errorf("Failed to resolve image, err: %s", e)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	ec.Log.Infof("Request to download thumbnail image: %+v", dimg)
	c.FileAttachment(dimg.Path, dimg.Name)
}

type TransferGalleryImageReq struct {
	Images []data.CreateGalleryImageCmd
}

/*
	Transfer image from file-server as a gallery image

	Request Body (JSON): TransferGalleryImageReq
*/
func TransferGalleryImageEndpoint(c *gin.Context, ec common.ExecContext) (any, error) {
	user := ec.User
	var cmd TransferGalleryImageReq
	server.MustBindJson(c, &cmd)

	if e := common.Validate(cmd); e != nil {
		return nil, e
	}

	if cmd.Images == nil || len(cmd.Images) < 1 {
		server.DispatchOk(c)
		return nil, nil
	}

	// validate the keys first
	for _, img := range cmd.Images {
		if isValid, e := client.ValidateFileKey(c.Request.Context(), img.FileKey, user.UserId); e != nil || !isValid {
			if e != nil {
				return nil, e
			}
			return nil, common.NewWebErr(fmt.Sprintf("Only file's owner can make it a gallery image ('%s')", img.Name))
		}
	}

	// start transferring
	go func(ec common.ExecContext, images []data.CreateGalleryImageCmd) {
		for _, cmd := range images {
			fi, er := client.GetFileInfo(ec.Ctx, cmd.FileKey)
			if er != nil {
				ec.Log.Errorf("Failed to fetch file info while transferring selected images, fi's fileKey: %s, error: %v", cmd.FileKey, er)
				continue
			}

			if fi.Data.FileType == client.FILE { // a file
				if data.GuessIsImage(fi.Data.Name, fi.Data.SizeInBytes) {
					nc := data.CreateGalleryImageCmd{GalleryNo: cmd.GalleryNo, Name: fi.Data.Name, FileKey: fi.Data.Uuid, FileLocalPath: fi.Data.LocalPath}
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

// Transfer image from file-server as a gallery image
// func TransferGalleryImageInDir(c *gin.Context, ec common.ExecContext) (any, error) {
// 	var cmd data.TransferGalleryImageInDirReq
// 	server.MustBindJson(c, &cmd)

// 	if e := common.Validate(cmd); e != nil {
// 		return nil, e
// 	}
// 	return nil, data.TransferImagesInDir(cmd, ec)
// }
