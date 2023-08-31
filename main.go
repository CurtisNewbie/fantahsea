package main

import (
	"os"

	"github.com/curtisnewbie/fantahsea/client"
	"github.com/curtisnewbie/fantahsea/data"
	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/goauth"
	"github.com/curtisnewbie/miso/bus"
	"github.com/curtisnewbie/miso/core"
	"github.com/curtisnewbie/miso/server"
	"github.com/gin-gonic/gin"
)

const (
	MNG_FILE_CODE = "manage-files"
	MNG_FILE_NAME = "Manage files"
)

func main() {

	server.PreServerBootstrap(func(c core.Rail) error {
		goauth.ReportResourcesOnBootstrapped(c, []goauth.AddResourceReq{{Code: MNG_FILE_CODE, Name: MNG_FILE_NAME}})
		goauth.ReportPathsOnBootstrapped(c)
		return nil
	})

	server.Get("/open/api/gallery/brief/owned",
		func(c *gin.Context, rail core.Rail) (any, error) {
			user := common.GetUser(rail)
			return data.ListOwnedGalleryBriefs(rail, user)
		},
		goauth.PathDocExtra(goauth.PathDoc{
			Desc: "List owned gallery brief info",
			Type: goauth.PT_PROTECTED,
			Code: MNG_FILE_CODE,
		}))

	server.IPost("/open/api/gallery/new",
		func(c *gin.Context, rail core.Rail, cmd data.CreateGalleryCmd) (*data.Gallery, error) {
			user := common.GetUser(rail)
			return data.CreateGallery(rail, cmd, user)
		},
		goauth.PathDocExtra(goauth.PathDoc{
			Desc: "Create new gallery",
			Type: goauth.PT_PROTECTED,
			Code: MNG_FILE_CODE,
		}))

	server.IPost("/open/api/gallery/update",
		func(c *gin.Context, rail core.Rail, cmd data.UpdateGalleryCmd) (any, error) {
			user := common.GetUser(rail)
			client.DispatchUserOpLog(rail, "UpdateGalleryEndpoint", "Update gallery", cmd, user)
			e := data.UpdateGallery(rail, cmd, user)
			return nil, e
		},
		goauth.PathDocExtra(goauth.PathDoc{
			Desc: "Update gallery",
			Type: goauth.PT_PROTECTED,
			Code: MNG_FILE_CODE,
		}))

	server.IPost("/open/api/gallery/delete",
		func(c *gin.Context, rail core.Rail, cmd data.DeleteGalleryCmd) (any, error) {
			user := common.GetUser(rail)
			client.DispatchUserOpLog(rail, "DeleteGalleryEndpoint", "Delete Gallery", cmd, user)
			e := data.DeleteGallery(rail, cmd, user)
			return nil, e
		},
		goauth.PathDocExtra(goauth.PathDoc{
			Desc: "Delete gallery",
			Type: goauth.PT_PROTECTED,
			Code: MNG_FILE_CODE,
		}))

	server.IPost("/open/api/gallery/list",
		func(c *gin.Context, rail core.Rail, cmd data.ListGalleriesCmd) (any, error) {
			user := common.GetUser(rail)
			return data.ListGalleries(rail, cmd, user)
		},
		goauth.PathDocExtra(goauth.PathDoc{
			Desc: "List galleries",
			Type: goauth.PT_PROTECTED,
			Code: MNG_FILE_CODE,
		}))

	server.IPost("/open/api/gallery/access/grant",
		func(c *gin.Context, rail core.Rail, cmd data.PermitGalleryAccessCmd) (any, error) {
			user := common.GetUser(rail)
			client.DispatchUserOpLog(rail, "GrantGalleryAccessEndpoint", "Grant access to the gallery", cmd, user)
			e := data.GrantGalleryAccessToUser(rail, cmd, user)
			return nil, e
		},
		goauth.PathDocExtra(goauth.PathDoc{
			Desc: "List granted access to the galleries",
			Type: goauth.PT_PROTECTED,
			Code: MNG_FILE_CODE,
		}))

	server.IPost("/open/api/gallery/images",
		func(c *gin.Context, rail core.Rail, cmd data.ListGalleryImagesCmd) (*data.ListGalleryImagesResp, error) {
			user := common.GetUser(rail)
			return data.ListGalleryImages(rail, cmd, user)
		},
		goauth.PathDocExtra(goauth.PathDoc{
			Desc: "List images of gallery",
			Type: goauth.PT_PROTECTED,
			Code: MNG_FILE_CODE,
		}))

	server.IPost("/open/api/gallery/image/transfer",
		func(c *gin.Context, rail core.Rail, cmd data.TransferGalleryImageReq) (any, error) {
			user := common.GetUser(rail)
			return data.BatchTransferAsync(rail, cmd, user)
		},
		goauth.PathDocExtra(goauth.PathDoc{
			Desc: "Host selected images on gallery",
			Type: goauth.PT_PROTECTED,
			Code: MNG_FILE_CODE,
		}))

	bus.DeclareEventBus(data.AddDirGalleryImageEventBus)
	bus.DeclareEventBus(data.NotifyFileDeletedEventBus)

	bus.SubscribeEventBus(data.AddDirGalleryImageEventBus, 2, func(rail core.Rail, evt data.CreateGalleryImgEvent) error {
		return data.OnCreateGalleryImgEvent(rail, evt)
	})

	bus.SubscribeEventBus(data.NotifyFileDeletedEventBus, 2, func(rail core.Rail, evt data.NotifyFileDeletedEvent) error {
		return data.OnNotifyFileDeletedEvent(rail, evt)
	})

	server.BootstrapServer(os.Args)
}
