package controller

import (
	"github.com/curtisnewbie/fantahsea/client"
	"github.com/curtisnewbie/fantahsea/data"
	"github.com/curtisnewbie/gocommon/common"
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
func ListGalleriesEndpoint(c *gin.Context, ec common.ExecContext, cmd data.ListGalleriesCmd) (any, error) {
	if e := common.Validate(cmd); e != nil {
		return nil, e
	}
	return data.ListGalleries(cmd, ec)
}

/*
	CreateGalleryEndpoint web endpoint

	Request Body (JSON): CreateGalleryCmd
*/
func CreateGalleryEndpoint(c *gin.Context, ec common.ExecContext, cmd data.CreateGalleryCmd) (any, error) {
	if e := common.Validate(cmd); e != nil {
		return nil, e
	}

	return data.CreateGallery(cmd, ec)
}

/*
	Update Gallery web endpoint

	Request Body (JSON): UpdateGalleryCmd
*/
func UpdateGalleryEndpoint(c *gin.Context, ec common.ExecContext, cmd data.UpdateGalleryCmd) (any, error) {
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
func DeleteGalleryEndpoint(c *gin.Context, ec common.ExecContext, cmd data.DeleteGalleryCmd) (any, error) {
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
func GrantGalleryAccessEndpoint(c *gin.Context, ec common.ExecContext, cmd data.PermitGalleryAccessCmd) (any, error) {
	client.DispatchUserOpLog(ec, "GrantGalleryAccessEndpoint", "Grant access to the gallery", cmd)

	if e := common.Validate(cmd); e != nil {
		return nil, e
	}

	if e := data.GrantGalleryAccessToUser(cmd, ec); e != nil {
		return nil, e
	}

	return nil, nil
}
