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
	router.POST(server.ResolvePath("/gallery/images", true), ListImagesEndpoint)
	router.GET(server.ResolvePath("/gallery/image/download", true), DownloadImageEndpoint)
	router.POST(server.ResolvePath("/gallery/image/transfer", true), TransferGalleryImageEndpoint)
}

/*
	List images of a gallery

	Request Body (JSON): ListGalleryImagesCmd
*/
func ListImagesEndpoint(c *gin.Context) {

	user, e := util.ExtractUser(c)
	if e != nil {
		util.DispatchErrJson(c, e)
		return
	}

	var cmd data.ListGalleryImagesCmd
	e = c.ShouldBindJSON(&cmd)
	if e != nil {
		util.DispatchErrJson(c, e)
		return
	}

	resp, e := data.ListGalleryImages(&cmd, user)
	if e != nil {
		util.DispatchErrJson(c, e)
		return
	}

	util.DispatchOkWData(c, resp)
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
func TransferGalleryImageEndpoint(c *gin.Context) {

	var user *util.User
	var e error

	if user, e = util.ExtractUser(c); e != nil {
		util.DispatchErrJson(c, e)
		return
	}

	var req TransferGalleryImageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		util.DispatchErrJson(c, err)
		return
	}

	if req.Images == nil {
		util.DispatchOk(c)
		return
	}

	// validate the keys first
	for _, cmd := range req.Images {
		if isValid, e := client.ValidateFileKey(cmd.FileKey, user.UserId); e != nil || !isValid {
			if e != nil {
				util.DispatchErrJson(c, e)
				return
			}
			util.DispatchErrJson(c, weberr.NewWebErr(fmt.Sprintf("Only file's owner can make it a gallery image ('%s')", cmd.Name)))
			return
		}
	}

	// start transferring
	go func(images []data.CreateGalleryImageCmd) {
		for _, cmd := range images {
			// todo Add a redis-lock for this method :D
			if e = data.CreateGalleryImage(&cmd, user); e != nil {
				log.Printf("Failed to transfer gallery image, e: %v", e)
				return
			}
		}
	}(req.Images)

	util.DispatchOk(c)
}
