package controller

import (
	"github.com/curtisnewbie/fantahsea/data"
	"github.com/curtisnewbie/gocommon/util"
	"github.com/curtisnewbie/gocommon/web/server"
	"github.com/gin-gonic/gin"
)

func RegisterGalleryRoutes(router *gin.Engine) {
	router.GET(server.ResolvePath("/gallery/brief/owned", true), util.BuildAuthRouteHandler(ListOwnedGalleryBriefsEndpoint))
	router.POST(server.ResolvePath("/gallery/new", true), util.BuildAuthRouteHandler(CreateGalleryEndpoint))
	router.POST(server.ResolvePath("/gallery/update", true), util.BuildAuthRouteHandler(UpdateGalleryEndpoint))
	router.POST(server.ResolvePath("/gallery/delete", true), util.BuildAuthRouteHandler(DeleteGalleryEndpoint))
	router.POST(server.ResolvePath("/gallery/list", true), util.BuildAuthRouteHandler(ListGalleriesEndpoint))
	router.POST(server.ResolvePath("/gallery/access/grant", true), util.BuildAuthRouteHandler(GrantGalleryAccessEndpoint))
}

// List owned gallery briefs list endpoint
func ListOwnedGalleryBriefsEndpoint(c *gin.Context, user *util.User) (any, error) {
	return data.ListOwnedGalleryBriefs(user)
}

/*
	ListGalleriesEndpoint web endpoint

	Request Body (JSON): ListGalleriesCmd
*/
func ListGalleriesEndpoint(c *gin.Context, user *util.User) (any, error) {
	var cmd data.ListGalleriesCmd
	util.MustBindJson(c, &cmd)

	return data.ListGalleries(&cmd, user)
}

/*
	CreateGalleryEndpoint web endpoint

	Request Body (JSON): CreateGalleryCmd
*/
func CreateGalleryEndpoint(c *gin.Context, user *util.User) (any, error) {
	var cmd data.CreateGalleryCmd
	util.MustBindJson(c, &cmd)

	result, er := util.LockRun("fantahsea:gallery:create:"+user.UserNo, func() any {
		if _, e := data.CreateGallery(&cmd, user); e != nil {
			return e
		}
		return nil
	})

	if er != nil {
		return nil, er
	}

	if result != nil {
		if casted, isOk := result.(error); isOk {
			return nil, casted
		}
	}
	return nil, nil
}

/*
	Update Gallery web endpoint

	Request Body (JSON): UpdateGalleryCmd
*/
func UpdateGalleryEndpoint(c *gin.Context, user *util.User) (any, error) {
	var cmd data.UpdateGalleryCmd
	util.MustBindJson(c, &cmd)

	if e := data.UpdateGallery(&cmd, user); e != nil {
		return nil, e
	}
	return nil, nil
}

// todo how about the temporary files we uploaded :D, need to handle them properly
/*
	Delete Gallery web endpoint

	Request Body (JSON): DeleteGalleryCmd
*/
func DeleteGalleryEndpoint(c *gin.Context, user *util.User) (any, error) {
	var cmd data.DeleteGalleryCmd
	util.MustBindJson(c, &cmd)

	if e := data.DeleteGallery(&cmd, user); e != nil {
		return nil, e
	}

	return nil, nil
}

/*
	Permit a user access to the gallery

	Request Body (JSON): PermitGalleryAccessCmd
*/
func GrantGalleryAccessEndpoint(c *gin.Context, user *util.User) (any, error) {
	var cmd data.PermitGalleryAccessCmd
	util.MustBindJson(c, &cmd)

	if e := data.GrantGalleryAccessToUser(&cmd, user); e != nil {
		return nil, e
	}

	return nil, nil
}
