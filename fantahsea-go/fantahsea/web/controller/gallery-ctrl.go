package controller

import (
	"fantahsea/data"
	"fantahsea/util"
	"fantahsea/web/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterGalleryRoutes(router *gin.Engine) {
	router.PUT(ResolvePath("/gallery/new", true), CreateGalleryEndpoint)
	router.POST(ResolvePath("/gallery/update", true), UpdateGalleryEndpoint)
	router.POST(ResolvePath("/gallery/delete", true), DeleteGalleryEndpoint)
	router.POST(ResolvePath("/gallery/list", true), ListGalleriesEndpoint)
	router.POST(ResolvePath("/gallery/access/grant", true), GrantGalleryAccessEndpoint)
}

/* ListGalleriesEndpoint web endpoint */
func ListGalleriesEndpoint(c *gin.Context) {
	user, err := util.ExtractUser(c)
	if err != nil {
		util.DispatchErrJson(c, err)
		return
	}

	var cmd data.ListGalleriesCmd
	if err := c.ShouldBindJSON(&cmd); err != nil {
		util.DispatchErrJson(c, err)
		return
	}

	resp, err := data.ListGalleries(&cmd, user)
	if err != nil {
		util.DispatchErrJson(c, err)
		return
	}

	util.DispatchOkWData(c, resp)
}

/* CreateGalleryEndpoint web endpoint */
func CreateGalleryEndpoint(c *gin.Context) {
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
func UpdateGalleryEndpoint(c *gin.Context) {
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
func DeleteGalleryEndpoint(c *gin.Context) {
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
func GrantGalleryAccessEndpoint(c *gin.Context) {
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

	if err := data.GrantGalleryAccessToUser(&cmd, user); err != nil {
		util.DispatchErrJson(c, err)
		return
	}

	util.DispatchOk(c)
}
