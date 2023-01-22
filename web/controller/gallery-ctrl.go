package controller

import (
	"github.com/curtisnewbie/fantahsea/client"
	"github.com/curtisnewbie/fantahsea/data"
	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/redis"
	"github.com/curtisnewbie/gocommon/server"
	"github.com/gin-gonic/gin"
)

// List owned gallery briefs list endpoint
func ListOwnedGalleryBriefsEndpoint(c *gin.Context, ec common.ExecContext) (any, error) {
	return data.ListOwnedGalleryBriefs(ec)
}

/*
	ListGalleriesEndpoint web endpoint

	Request Body (JSON): ListGalleriesCmd
*/
func ListGalleriesEndpoint(c *gin.Context, ec common.ExecContext) (any, error) {
	var cmd data.ListGalleriesCmd
	server.MustBindJson(c, &cmd)
	if e := common.Validate(cmd); e != nil {
		return nil, e
	}
	return data.ListGalleries(cmd, ec)
}

/*
	CreateGalleryEndpoint web endpoint

	Request Body (JSON): CreateGalleryCmd
*/
func CreateGalleryEndpoint(c *gin.Context, ec common.ExecContext) (any, error) {
	var cmd data.CreateGalleryCmd
	server.MustBindJson(c, &cmd)

	if e := common.Validate(cmd); e != nil {
		return nil, e
	}

	_, er := redis.RLockRun("fantahsea:gallery:create:"+ec.User.UserNo, func() (any, error) {
		return data.CreateGallery(cmd, ec)
	})

	if er != nil {
		return nil, er
	}

	return nil, nil
}

/*
	Update Gallery web endpoint

	Request Body (JSON): UpdateGalleryCmd
*/
func UpdateGalleryEndpoint(c *gin.Context, ec common.ExecContext) (any, error) {
	var cmd data.UpdateGalleryCmd
	server.MustBindJson(c, &cmd)

	client.DispatchUserOpLog(ec, "UpdateGalleryEndpoint", "Update gallery", cmd)

	if e := common.Validate(cmd); e != nil {
		return nil, e
	}

	if e := data.UpdateGallery(cmd, ec); e != nil {
		return nil, e
	}
	return nil, nil
}

/*
	Delete Gallery web endpoint

	Request Body (JSON): DeleteGalleryCmd
*/
func DeleteGalleryEndpoint(c *gin.Context, ec common.ExecContext) (any, error) {
	var cmd data.DeleteGalleryCmd
	server.MustBindJson(c, &cmd)

	client.DispatchUserOpLog(ec, "DeleteGalleryEndpoint", "Delete Gallery", cmd)

	if e := common.Validate(cmd); e != nil {
		return nil, e
	}

	if e := data.DeleteGallery(cmd, ec); e != nil {
		return nil, e
	}

	return nil, nil
}

/*
	Permit a user access to the gallery

	Request Body (JSON): PermitGalleryAccessCmd
*/
func GrantGalleryAccessEndpoint(c *gin.Context, ec common.ExecContext) (any, error) {
	var cmd data.PermitGalleryAccessCmd
	server.MustBindJson(c, &cmd)

	client.DispatchUserOpLog(ec, "GrantGalleryAccessEndpoint", "Grant access to the gallery", cmd)

	if e := common.Validate(cmd); e != nil {
		return nil, e
	}

	if e := data.GrantGalleryAccessToUser(cmd, ec); e != nil {
		return nil, e
	}

	return nil, nil
}
