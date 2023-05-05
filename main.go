package main

import (
	"log"
	"os"

	"github.com/curtisnewbie/fantahsea/data"
	"github.com/curtisnewbie/fantahsea/web/controller"
	"github.com/curtisnewbie/goauth/client/goauth-client-go/gclient"
	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/server"
)

const (
	MNG_FILE_CODE = "manage-files"
	MNG_FILE_NAME = "Manage files"
)

func main() {
	// register jobs
	common.ScheduleCron("0 0/10 * * * *", data.CleanUpDeletedGallery)
	ec := common.EmptyExecContext()

	if gclient.IsEnabled() {
		server.OnServerBootstrapped(func() {
			if e := gclient.AddResource(ec.Ctx, gclient.AddResourceReq{Code: MNG_FILE_CODE, Name: MNG_FILE_NAME}); e != nil {
				log.Fatalf("gclient.AddResource, %v", e)
			}
		})

		gclient.ReportPathsOnBootstrapped()
	}

	// public routes
	server.RawGet(server.OpenApiPath("/gallery/image/download"), controller.DownloadImageThumbnailEndpoint,
		gclient.PathDocExtra(gclient.PathDoc{Type: gclient.PT_PUBLIC, Desc: "Download gallery image by token"}))

	// authenticated routes
	server.Get(server.OpenApiPath("/gallery/brief/owned"), controller.ListOwnedGalleryBriefsEndpoint,
		gclient.PathDocExtra(gclient.PathDoc{Type: gclient.PT_PROTECTED, Desc: "List owned gallery brief info"}))
	server.PostJ(server.OpenApiPath("/gallery/new"), controller.CreateGalleryEndpoint,
		gclient.PathDocExtra(gclient.PathDoc{Type: gclient.PT_PROTECTED, Desc: "Create new gallery"}))
	server.PostJ(server.OpenApiPath("/gallery/update"), controller.UpdateGalleryEndpoint,
		gclient.PathDocExtra(gclient.PathDoc{Type: gclient.PT_PROTECTED, Desc: "Update gallery"}))
	server.PostJ(server.OpenApiPath("/gallery/delete"), controller.DeleteGalleryEndpoint,
		gclient.PathDocExtra(gclient.PathDoc{Type: gclient.PT_PROTECTED, Desc: "Delete gallery"}))
	server.PostJ(server.OpenApiPath("/gallery/list"), controller.ListGalleriesEndpoint,
		gclient.PathDocExtra(gclient.PathDoc{Type: gclient.PT_PROTECTED, Desc: "List galleries"}))
	server.PostJ(server.OpenApiPath("/gallery/access/grant"), controller.GrantGalleryAccessEndpoint,
		gclient.PathDocExtra(gclient.PathDoc{Type: gclient.PT_PROTECTED, Desc: "List granted access to the galleries"}))
	server.PostJ(server.OpenApiPath("/gallery/images"), controller.ListImagesEndpoint,
		gclient.PathDocExtra(gclient.PathDoc{Type: gclient.PT_PROTECTED, Desc: "List images of gallery"}))
	server.PostJ(server.OpenApiPath("/gallery/image/transfer"), controller.TransferGalleryImageEndpoint,
		gclient.PathDocExtra(gclient.PathDoc{Type: gclient.PT_PROTECTED, Desc: "Host selected images on gallery"}))

	// bootstrap server
	server.DefaultBootstrapServer(os.Args)
}
