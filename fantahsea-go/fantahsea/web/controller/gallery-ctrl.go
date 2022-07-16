package controller

import (
	"fantahsea/data"
	"fantahsea/util"
	"fantahsea/web/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterGalleryRoutes(router *gin.Engine) {
	router.PUT(ResolvePath("/gallery/new"), CreateGallery)
	router.POST(ResolvePath("/gallery/update"), UpdateGallery)
	router.POST(ResolvePath("/gallery/delete"), DeleteGallery)
	router.POST(ResolvePath("/gallery/list"), ListGalleries)
	router.POST(ResolvePath("/gallery/access"), PermitGalleryAccess)
}

/* ListGalleries web endpoint */
func ListGalleries(c *gin.Context) {
	user, err := util.ExtractUser(c)
	if err != nil {
		util.DispatchErrJson(c, err)
		return
	}

	var page dto.Paging
	if err := c.ShouldBindJSON(&page); err != nil {
		util.DispatchErrJson(c, err)
		return
	}

	resp, err := data.ListGalleries(&page, user)
	if err != nil {
		util.DispatchErrJson(c, err)
		return
	}

	util.DispatchOkWData(c, resp)
}

/* CreateGallery web endpoint */
func CreateGallery(c *gin.Context) {
	user, err := util.ExtractUser(c)
	if err != nil {
		util.DispatchErrJson(c, err)
		return
	}

	var cmd data.CreateGalleryCmd
	if err := c.ShouldBindJSON(&cmd); err != nil {
		util.DispatchJson(c, dto.ErrorResp("Illegal Arguments"))
		return
	}

	if _, err := data.CreateGallery(&cmd, user); err != nil {
		util.DispatchErrJson(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.OkResp())
}

/* Update Gallery web endpoint */
func UpdateGallery(c *gin.Context) {
	user, err := util.ExtractUser(c)
	if err != nil {
		util.DispatchErrJson(c, err)
		return
	}

	var cmd data.UpdateGalleryCmd
	if err := c.ShouldBindJSON(&cmd); err != nil {
		util.DispatchJson(c, dto.ErrorResp("Illegal Arguments"))
		return
	}

	if err := data.UpdateGallery(&cmd, user); err != nil {
		util.DispatchErrJson(c, err)
		return
	}

	util.DispatchOk(c)
}

/* Delete Gallery web endpoint */
func DeleteGallery(c *gin.Context) {
	user, err := util.ExtractUser(c)
	if err != nil {
		util.DispatchErrJson(c, err)
		return
	}

	var cmd data.DeleteGalleryCmd
	if err := c.ShouldBindJSON(&cmd); err != nil {
		util.DispatchJson(c, dto.ErrorResp("Illegal Arguments"))
		return
	}

	if err := data.DeleteGallery(&cmd, user); err != nil {
		util.DispatchErrJson(c, err)
		return
	}

	util.DispatchOk(c)
}

/* Permit a user access to the gallery */
func PermitGalleryAccess(c *gin.Context) {
	user, err := util.ExtractUser(c)
	if err != nil {
		util.DispatchErrJson(c, err)
		return
	}

	var cmd data.PermitGalleryAccessCmd
	if err := c.ShouldBindJSON(&cmd); err != nil {
		util.DispatchJson(c, dto.ErrorResp("Illegal Arguments"))
		return
	}

	if err := data.PermitGalleryAccess(&cmd, user); err != nil {
		util.DispatchErrJson(c, err)
		return
	}

	util.DispatchOk(c)
}
