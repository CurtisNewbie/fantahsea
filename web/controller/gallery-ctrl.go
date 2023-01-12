package controller

import (
	"github.com/curtisnewbie/fantahsea/data"
	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/redis"
	"github.com/curtisnewbie/gocommon/server"
	"github.com/gin-gonic/gin"
)

// List owned gallery briefs list endpoint
func ListOwnedGalleryBriefsEndpoint(c *gin.Context, user *common.User) (any, error) {
	return data.ListOwnedGalleryBriefs(user)
}

/*
	ListGalleriesEndpoint web endpoint

	Request Body (JSON): ListGalleriesCmd
*/
func ListGalleriesEndpoint(c *gin.Context, user *common.User) (any, error) {
	var cmd data.ListGalleriesCmd
	server.MustBindJson(c, &cmd)

	return data.ListGalleries(&cmd, user)
}

/*
	CreateGalleryEndpoint web endpoint

	Request Body (JSON): CreateGalleryCmd
*/
func CreateGalleryEndpoint(c *gin.Context, user *common.User) (any, error) {
	var cmd data.CreateGalleryCmd
	server.MustBindJson(c, &cmd)

	result, er := redis.RLockRun("fantahsea:gallery:create:"+user.UserNo, func() any {
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
func UpdateGalleryEndpoint(c *gin.Context, user *common.User) (any, error) {
	var cmd data.UpdateGalleryCmd
	server.MustBindJson(c, &cmd)

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
func DeleteGalleryEndpoint(c *gin.Context, user *common.User) (any, error) {
	var cmd data.DeleteGalleryCmd
	server.MustBindJson(c, &cmd)

	if e := data.DeleteGallery(&cmd, user); e != nil {
		return nil, e
	}

	return nil, nil
}

/*
	Permit a user access to the gallery

	Request Body (JSON): PermitGalleryAccessCmd
*/
func GrantGalleryAccessEndpoint(c *gin.Context, user *common.User) (any, error) {
	var cmd data.PermitGalleryAccessCmd
	server.MustBindJson(c, &cmd)

	if e := data.GrantGalleryAccessToUser(&cmd, user); e != nil {
		return nil, e
	}

	return nil, nil
}
