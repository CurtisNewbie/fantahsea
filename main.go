package main

import (
	"os"

	"github.com/curtisnewbie/fantahsea/client"
	"github.com/curtisnewbie/fantahsea/data"
	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/goauth"
	"github.com/curtisnewbie/gocommon/server"
	"github.com/gin-gonic/gin"
)

const (
	MNG_FILE_CODE = "manage-files"
	MNG_FILE_NAME = "Manage files"

	CREATE_GALLERY_IMAGE_EVENT_BUS = "fantahsea:gallery:image:create"
)

func main() {

	server.PreServerBootstrap(func(c common.ExecContext) error {
		if goauth.IsEnabled() {
			server.PostServerBootstrapped(func(c common.ExecContext) error {
				return goauth.AddResource(c.Ctx, goauth.AddResourceReq{Code: MNG_FILE_CODE, Name: MNG_FILE_NAME})
			})

			goauth.ReportPathsOnBootstrapped()
		}
		return nil
	})

	server.Get("/open/api/gallery/brief/owned",
		func(c *gin.Context, ec common.ExecContext) (any, error) {
			return data.ListOwnedGalleryBriefs(ec)
		},
		goauth.PathDocExtra(goauth.PathDoc{
			Desc: "List owned gallery brief info",
			Type: goauth.PT_PROTECTED,
			Code: MNG_FILE_CODE,
		}))

	server.IPost("/open/api/gallery/new",
		func(c *gin.Context, ec common.ExecContext, cmd data.CreateGalleryCmd) (any, error) {
			return data.CreateGallery(cmd, ec)
		},
		goauth.PathDocExtra(goauth.PathDoc{
			Desc: "Create new gallery",
			Type: goauth.PT_PROTECTED,
			Code: MNG_FILE_CODE,
		}))

	server.IPost("/open/api/gallery/update",
		func(c *gin.Context, ec common.ExecContext, cmd data.UpdateGalleryCmd) (any, error) {
			client.DispatchUserOpLog(ec, "UpdateGalleryEndpoint", "Update gallery", cmd)
			e := data.UpdateGallery(cmd, ec)
			return nil, e
		},
		goauth.PathDocExtra(goauth.PathDoc{
			Desc: "Update gallery",
			Type: goauth.PT_PROTECTED,
			Code: MNG_FILE_CODE,
		}))

	server.IPost("/open/api/gallery/delete",
		func(c *gin.Context, ec common.ExecContext, cmd data.DeleteGalleryCmd) (any, error) {
			client.DispatchUserOpLog(ec, "DeleteGalleryEndpoint", "Delete Gallery", cmd)
			e := data.DeleteGallery(cmd, ec)
			return nil, e
		},
		goauth.PathDocExtra(goauth.PathDoc{
			Desc: "Delete gallery",
			Type: goauth.PT_PROTECTED,
			Code: MNG_FILE_CODE,
		}))

	server.IPost("/open/api/gallery/list",
		func(c *gin.Context, ec common.ExecContext, cmd data.ListGalleriesCmd) (any, error) {
			return data.ListGalleries(cmd, ec)
		},
		goauth.PathDocExtra(goauth.PathDoc{
			Desc: "List galleries",
			Type: goauth.PT_PROTECTED,
			Code: MNG_FILE_CODE,
		}))

	server.IPost("/open/api/gallery/access/grant",
		func(c *gin.Context, ec common.ExecContext, cmd data.PermitGalleryAccessCmd) (any, error) {
			client.DispatchUserOpLog(ec, "GrantGalleryAccessEndpoint", "Grant access to the gallery", cmd)
			e := data.GrantGalleryAccessToUser(cmd, ec)
			return nil, e
		},
		goauth.PathDocExtra(goauth.PathDoc{
			Desc: "List granted access to the galleries",
			Type: goauth.PT_PROTECTED,
			Code: MNG_FILE_CODE,
		}))

	server.IPost("/open/api/gallery/images",
		func(c *gin.Context, ec common.ExecContext, cmd data.ListGalleryImagesCmd) (*data.ListGalleryImagesResp, error) {
			return data.ListGalleryImages(cmd, ec)
		},
		goauth.PathDocExtra(goauth.PathDoc{
			Desc: "List images of gallery",
			Type: goauth.PT_PROTECTED,
			Code: MNG_FILE_CODE,
		}))

	server.IPost("/open/api/gallery/image/transfer",
		func(c *gin.Context, ec common.ExecContext, cmd data.TransferGalleryImageReq) (any, error) {
			return data.BatchTransferAsync(ec, cmd)
		},
		goauth.PathDocExtra(goauth.PathDoc{
			Desc: "Host selected images on gallery",
			Type: goauth.PT_PROTECTED,
			Code: MNG_FILE_CODE,
		}))

	server.BootstrapServer(os.Args)
}
