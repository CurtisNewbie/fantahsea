package controller

import (
	"fantahsea/client"
	. "fantahsea/data"
	"fantahsea/err"
	. "fantahsea/util"
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

	user, e := ExtractUser(c)
	if e != nil {
		DispatchErrJson(c, e)
		return
	}

	var cmd ListGalleryImagesCmd
	e = c.ShouldBindJSON(&cmd)
	if e != nil {
		DispatchErrJson(c, e)
		return
	}

	resp, e := ListGalleryImages(&cmd, user)
	if e != nil {
		DispatchErrJson(c, e)
		return
	}

	DispatchOkWData(c, resp)
}

/*
	Download image

	Query Param: imageNo
*/
func DownloadImageEndpoint(c *gin.Context) {
	user, e := ExtractUser(c)
	if e != nil {
		DispatchErrJson(c, e)
		return
	}

	dimg, e := ResolveImageDInfo(c.Query("imageNo"), c.Query("thumbnail"), user)
	if e != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	log.Infof("Request to download gallery image: %+v", dimg)

	// Write the file to client
	c.FileAttachment(dimg.Path, dimg.Name)
}

type TransferGalleryImageReq struct {
	Images []CreateGalleryImageCmd
}

/*
	Transfer image from file-server to fantahsea as a gallery image

	Request Body (JSON): TransferGalleryImageReq
*/
func TransferGalleryImageEndpoint(c *gin.Context) {

	var user *User
	var e error

	if user, e = ExtractUser(c); e != nil {
		DispatchErrJson(c, e)
		return
	}

	var req TransferGalleryImageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		DispatchErrJson(c, err)
		return
	}

	for _, cmd := range req.Images {

		// validate the key first
		if isValid, e := client.ValidateFileKey(cmd.FileKey, user.UserId); e != nil || !isValid {
			if e != nil {
				DispatchErrJson(c, e)
				return
			}
			DispatchErrJson(c, err.NewWebErr(fmt.Sprintf("Only file's owner can make it a gallery image ('%s')", cmd.Name)))
			return
		}

		// create record
		if e = CreateGalleryImage(&cmd, user); e != nil {
			DispatchErrJson(c, e)
			return
		}

	}

	DispatchOk(c)
}
