package controller

import (
	"fmt"
	"net/http"

	"github.com/curtisnewbie/fantahsea/data"
	"github.com/curtisnewbie/fantahsea/fclient"
	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/server"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

/*
	List images of a gallery

	Request Body (JSON): ListGalleryImagesCmd
*/
func ListImagesEndpoint(c *gin.Context, user *common.User) (any, error) {
	var cmd data.ListGalleryImagesCmd
	server.MustBindJson(c, &cmd)

	return data.ListGalleryImages(&cmd, user)
}

/*
	Download image

	Query Param: imageNo
*/
func DownloadImageEndpoint(c *gin.Context) {

	token, thumbnail := c.Query("token"), c.Query("thumbnail")

	log.Printf("Download Image, token: %s, thumbnail: %s", token, thumbnail)
	dimg, e := data.ResolveImageDInfo(token, thumbnail)
	if e != nil {
		log.Printf("Failed to resolve image, err: %s", e)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	log.Infof("Request to download gallery image: %+v", dimg)

	// Write the file to client
	c.FileAttachment(dimg.Path, dimg.Name)
}

type TransferGalleryImageReq struct {
	Images []data.CreateGalleryImageCmd
}

/*
	Transfer image from file-server as a gallery image

	Request Body (JSON): TransferGalleryImageReq
*/
func TransferGalleryImageEndpoint(c *gin.Context, user *common.User) (any, error) {

	var req TransferGalleryImageReq
	server.MustBindJson(c, &req)

	if req.Images == nil {
		server.DispatchOk(c)
		return nil, nil
	}

	// validate the keys first
	for _, cmd := range req.Images {
		if isValid, e := fclient.ValidateFileKey(c.Request.Context(), cmd.FileKey, user.UserId); e != nil || !isValid {
			if e != nil {
				return nil, e
			}
			return nil, common.NewWebErr(fmt.Sprintf("Only file's owner can make it a gallery image ('%s')", cmd.Name))
		}
	}

	// start transferring
	go func(images []data.CreateGalleryImageCmd) {
		for _, cmd := range images {
			// todo Add a redis-lock for this method :D, a unique constraint for gallery_no & file_key should do for now
			if e := data.CreateGalleryImage(c.Request.Context(), &cmd, user); e != nil {
				log.Printf("Failed to transfer gallery image, e: %v", e)
				return
			}
		}
	}(req.Images)

	return nil, nil
}

// Transfer image from file-server as a gallery image
func TransferGalleryImageInDir(c *gin.Context, user *common.User) (any, error) {
	var req data.TransferGalleryImageInDirReq
	server.MustBindJson(c, &req)
	e := data.TransferImagesInDir(c.Request.Context(), &req, user)
	return nil, e
}
