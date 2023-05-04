package main

import (
	"log"
	"os"

	"github.com/curtisnewbie/fantahsea/data"
	"github.com/curtisnewbie/fantahsea/web/controller"
	"github.com/curtisnewbie/goauth/client/goauth-client-go/gclient"
	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/server"
	"github.com/gin-gonic/gin"
)

const (
	MNG_FILE_CODE = "manage-files"
	MNG_FILE_NAME = "Manage files"
)

func main() {
	// register jobs
	common.ScheduleCron("0 0/10 * * * *", data.CleanUpDeletedGallery)
	ec := common.EmptyExecContext()

	server.OnServerBootstrapped(func() {
		if e := gclient.AddResource(ec.Ctx, gclient.AddResourceReq{Code: MNG_FILE_CODE, Name: MNG_FILE_NAME}); e != nil {
			log.Fatalf("gclient.AddResource, %v", e)
		}
	})

	// public routes
	gclient.RawGet(server.OpenApiPath("/gallery/image/download"), func(c *gin.Context, ec common.ExecContext) {
		controller.DownloadImageThumbnailEndpoint(c, ec)
	}, gclient.PathDoc{Type: gclient.PT_PUBLIC, Desc: "Download gallery image by token"})

	// authenticated routes
	gclient.Get(server.OpenApiPath("/gallery/brief/owned"), controller.ListOwnedGalleryBriefsEndpoint, gclient.PathDoc{Type: gclient.PT_PROTECTED,
		Desc: "List owned gallery brief info"})

	gclient.PostJ(server.OpenApiPath("/gallery/new"), controller.CreateGalleryEndpoint, gclient.PathDoc{
		Type: gclient.PT_PROTECTED, Desc: "Create new gallery",
	})
	gclient.PostJ(server.OpenApiPath("/gallery/update"), controller.UpdateGalleryEndpoint, gclient.PathDoc{
		Type: gclient.PT_PROTECTED, Desc: "Update gallery",
	})

	gclient.PostJ(server.OpenApiPath("/gallery/delete"), controller.DeleteGalleryEndpoint, gclient.PathDoc{
		Type: gclient.PT_PROTECTED, Desc: "Delete gallery",
	})

	gclient.PostJ(server.OpenApiPath("/gallery/list"), controller.ListGalleriesEndpoint, gclient.PathDoc{
		Type: gclient.PT_PROTECTED, Desc: "List galleries",
	})

	gclient.PostJ(server.OpenApiPath("/gallery/access/grant"), controller.GrantGalleryAccessEndpoint, gclient.PathDoc{
		Type: gclient.PT_PROTECTED, Desc: "List granted access to the galleries",
	})

	gclient.PostJ(server.OpenApiPath("/gallery/images"), controller.ListImagesEndpoint, gclient.PathDoc{
		Type: gclient.PT_PROTECTED, Desc: "List images of gallery",
	})

	gclient.PostJ(server.OpenApiPath("/gallery/image/transfer"), controller.TransferGalleryImageEndpoint, gclient.PathDoc{
		Type: gclient.PT_PROTECTED, Desc: "Host selected images on gallery",
	})

	// bootstrap server
	server.DefaultBootstrapServer(os.Args)
}
