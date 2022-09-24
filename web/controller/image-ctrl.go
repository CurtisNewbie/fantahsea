package controller

import (
	"fmt"
	"net/http"

	"github.com/curtisnewbie/fantahsea/client"
	"github.com/curtisnewbie/fantahsea/data"
	"github.com/curtisnewbie/gocommon/util"
	"github.com/curtisnewbie/gocommon/web/server"
	"github.com/curtisnewbie/gocommon/weberr"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Register routes
func RegisterGalleryImageRoutes(router *gin.Engine) {
	router.POST(server.ResolvePath("/gallery/images", true), util.BuildAuthRouteHandler(ListImagesEndpoint))
	router.GET(server.ResolvePath("/gallery/image/download", true), DownloadImageEndpoint)
	router.POST(server.ResolvePath("/gallery/image/transfer", true), util.BuildAuthRouteHandler(TransferGalleryImageEndpoint))
}

/*
	List images of a gallery

	Request Body (JSON): ListGalleryImagesCmd
*/
func ListImagesEndpoint(c *gin.Context, user *util.User) (any, error) {
	var cmd data.ListGalleryImagesCmd
	util.MustBindJson(c, &cmd)

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
func TransferGalleryImageEndpoint(c *gin.Context, user *util.User) (any, error) {

	var req TransferGalleryImageReq
	util.MustBindJson(c, &req)

	if req.Images == nil {
		util.DispatchOk(c)
		return nil, nil
	}

	// validate the keys first
	for _, cmd := range req.Images {
		if isValid, e := client.ValidateFileKey(cmd.FileKey, user.UserId); e != nil || !isValid {
			if e != nil {
				return nil, e
			}
			return nil, weberr.NewWebErr(fmt.Sprintf("Only file's owner can make it a gallery image ('%s')", cmd.Name))
		}
	}

	// start transferring
	go func(images []data.CreateGalleryImageCmd) {
		for _, cmd := range images {
			// todo Add a redis-lock for this method :D, a unique constraint for gallery_no & file_key should do for now
			if e := data.CreateGalleryImage(&cmd, user); e != nil {
				log.Printf("Failed to transfer gallery image, e: %v", e)
				return
			}
		}
	}(req.Images)

	return nil, nil
}
