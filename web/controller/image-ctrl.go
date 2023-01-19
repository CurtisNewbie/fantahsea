package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/curtisnewbie/fantahsea/data"
	"github.com/curtisnewbie/fantahsea/fclient"
	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/server"
	"github.com/gin-gonic/gin"
)

/*
	List images of a gallery

	Request Body (JSON): ListGalleryImagesCmd
*/
func ListImagesEndpoint(c *gin.Context, ec server.ExecContext) (any, error) {
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
func DownloadImageThumbnailEndpoint(c *gin.Context, ec server.ExecContext) {
	token := c.Query("token")
	ec.Log.Printf("Download Image thumbnail, token: %s", token)
	dimg, e := data.ResolveImageThumbnail(token)
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
func TransferGalleryImageEndpoint(c *gin.Context, ec server.ExecContext) (any, error) {
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
		if isValid, e := fclient.ValidateFileKey(c.Request.Context(), img.FileKey, user.UserId); e != nil || !isValid {
			if e != nil {
				return nil, e
			}
			return nil, common.NewWebErr(fmt.Sprintf("Only file's owner can make it a gallery image ('%s')", img.Name))
		}
	}

	// start transferring
	go func(req server.ExecContext, images []data.CreateGalleryImageCmd) {
		for _, cmd := range images {
			// todo Add a redis-lock for this method :D, a unique constraint for gallery_no & file_key should do for now
			if e := data.CreateGalleryImage(req, cmd); e != nil {
				log.Printf("Failed to transfer gallery image, e: %v", e)
				return
			}
		}
	}(ec, cmd.Images)

	return nil, nil
}

// Transfer image from file-server as a gallery image
func TransferGalleryImageInDir(c *gin.Context, ec server.ExecContext) (any, error) {
	var cmd data.TransferGalleryImageInDirReq
	server.MustBindJson(c, &cmd)

	if e := common.Validate(cmd); e != nil {
		return nil, e
	}
	return nil, data.TransferImagesInDir(cmd, ec)
}
