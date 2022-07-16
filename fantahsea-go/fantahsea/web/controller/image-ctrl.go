package controller

import (
	"fantahsea/client"
	"fantahsea/config"
	. "fantahsea/data"
	"fantahsea/err"
	. "fantahsea/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterGalleryImageRoutes(router *gin.Engine) {
	router.POST(ResolvePath("/gallery/images", true), ListImagesEndpoint)
	router.GET(ResolvePath("/gallery/image:imageNo", true), DownloadImageEndpoint)
	router.POST(ResolvePath("/gallery/image/transfer", true), TransferGalleryImageEndpoint)
}

// List images of a gallery
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

// Download image
func DownloadImageEndpoint(c *gin.Context) {
	user, e := ExtractUser(c)
	if e != nil {
		DispatchErrJson(c, e)
		return
	}

	dimg, e := ResolveImageDInfo(c.Param("imageNo"), user)
	if e != nil {
		DispatchErrJson(c, e)
		return
	}

	// Write the file to client
	c.FileAttachment(dimg.Path, dimg.Name)
}

/*
	Transfer image from file-server to fantahsea as a gallery image
*/
func TransferGalleryImageEndpoint(c *gin.Context) {

	var user *User
	var e error

	if user, e = ExtractUser(c); e != nil {
		DispatchErrJson(c, e)
		return
	}

	var cmd CreateGalleryImageCmd
	if err := c.ShouldBindJSON(&cmd); err != nil {
		DispatchErrJson(c, err)
		return
	}

	// validate the key first
	if isValid, e := client.ValidateFileKey(cmd.FileKey, user); e != nil || !isValid {
		if e != nil {
			DispatchErrJson(c, e)
			return
		}
		DispatchErrJson(c, err.NewWebErr("Only file's owner can make it a gallery image"))
		return
	}

	te := config.GetDB().Transaction(func(tx *gorm.DB) error {

		// create record
		var resp *CreateGalleryImageResp
		if resp, e = CreateGalleryImage(&cmd, user); e != nil {
			return e
		}

		// download the file from file-server
		if e := client.DownloadFile(resp.FileKey, user, resp.AbsPath); e != nil {
			return e
		}
		return nil
	})

	if te != nil {
		DispatchErrJson(c, e)
		return
	}

	DispatchOk(c)
}
