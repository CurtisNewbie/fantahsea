package controller

import (
	"fantahsea/client"
	"fantahsea/data"
	"fantahsea/util"
	"fantahsea/weberr"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Register routes
func RegisterGalleryImageRoutes(router *gin.Engine) {
	router.POST(ResolvePath("/gallery/images", true), ListImagesEndpoint)
	router.GET(ResolvePath("/gallery/image/download", true), DownloadImageEndpoint)
	router.POST(ResolvePath("/gallery/image/transfer", true), TransferGalleryImageEndpoint)
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
	Transfer image from file-server to fantahsea as a gallery image

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

	count := len(req.Images)
	for _, cmd := range req.Images {

		// validate the key first
		if isValid, e := client.ValidateFileKey(cmd.FileKey, user.UserId); e != nil || !isValid {
			if e != nil {
				util.DispatchErrJson(c, e)
				return
			}
			util.DispatchErrJson(c, weberr.NewWebErr(fmt.Sprintf("Only file's owner can make it a gallery image ('%s')", cmd.Name)))
			return
		}

		if e = data.CreateGalleryImage(&cmd, user); e != nil {
			if count < 2 {
				util.DispatchErrJson(c, e)
				return
			}
			log.Printf("Failed to transfer gallery image, e: %v", e)
		}

	}

	util.DispatchOk(c)
}
