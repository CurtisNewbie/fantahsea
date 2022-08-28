package controller

import (
	"github.com/curtisnewbie/fantahsea/data"
	"github.com/curtisnewbie/gocommon/util"
	"github.com/curtisnewbie/gocommon/web/dto"
	"github.com/curtisnewbie/gocommon/web/server"
	"github.com/gin-gonic/gin"
)

func RegisterGalleryRoutes(router *gin.Engine) {
	router.POST(server.ResolvePath("/gallery/new", true), CreateGalleryEndpoint)
	router.POST(server.ResolvePath("/gallery/update", true), UpdateGalleryEndpoint)
	router.POST(server.ResolvePath("/gallery/delete", true), DeleteGalleryEndpoint)
	router.POST(server.ResolvePath("/gallery/list", true), ListGalleriesEndpoint)
	router.POST(server.ResolvePath("/gallery/access/grant", true), GrantGalleryAccessEndpoint)
}

/*
	ListGalleriesEndpoint web endpoint

	Request Body (JSON): ListGalleriesCmd
*/
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

/*
	CreateGalleryEndpoint web endpoint

	Request Body (JSON): CreateGalleryCmd
*/
func CreateGalleryEndpoint(c *gin.Context) {
	user, e := util.ExtractUser(c)
	if e != nil {
		util.DispatchErrJson(c, e)
		return
	}

	var cmd data.CreateGalleryCmd
	if e := c.ShouldBindJSON(&cmd); e != nil {
		util.DispatchJson(c, dto.ErrorResp("Illegal Arguments"))
		return
	}

	if _, e := data.CreateGallery(&cmd, user); e != nil {
		util.DispatchErrJson(c, e)
		return
	}

	util.DispatchOk(c)
}

/*
	Update Gallery web endpoint

	Request Body (JSON): UpdateGalleryCmd
*/
func UpdateGalleryEndpoint(c *gin.Context) {
	user, e := util.ExtractUser(c)
	if e != nil {
		util.DispatchErrJson(c, e)
		return
	}

	var cmd data.UpdateGalleryCmd
	if e := c.ShouldBindJSON(&cmd); e != nil {
		util.DispatchJson(c, dto.ErrorResp("Illegal Arguments"))
		return
	}

	if e := data.UpdateGallery(&cmd, user); e != nil {
		util.DispatchErrJson(c, e)
		return
	}

	util.DispatchOk(c)
}

// todo how about the temporary files we uploaded :D, need to handle them properly
/*
	Delete Gallery web endpoint

	Request Body (JSON): DeleteGalleryCmd
*/
func DeleteGalleryEndpoint(c *gin.Context) {
	user, e := util.ExtractUser(c)
	if e != nil {
		util.DispatchErrJson(c, e)
		return
	}

	var cmd data.DeleteGalleryCmd
	if e := c.ShouldBindJSON(&cmd); e != nil {
		util.DispatchJson(c, dto.ErrorResp("Illegal Arguments"))
		return
	}

	if e := data.DeleteGallery(&cmd, user); e != nil {
		util.DispatchErrJson(c, e)
		return
	}

	util.DispatchOk(c)
}

/*
	Permit a user access to the gallery

	Request Body (JSON): PermitGalleryAccessCmd
*/
func GrantGalleryAccessEndpoint(c *gin.Context) {
	user, e := util.ExtractUser(c)
	if e != nil {
		util.DispatchErrJson(c, e)
		return
	}

	var cmd data.PermitGalleryAccessCmd
	if e := c.ShouldBindJSON(&cmd); e != nil {
		util.DispatchJson(c, dto.ErrorResp("Illegal Arguments"))
		return
	}

	if e := data.GrantGalleryAccessToUser(&cmd, user); e != nil {
		util.DispatchErrJson(c, e)
		return
	}

	util.DispatchOk(c)
}
