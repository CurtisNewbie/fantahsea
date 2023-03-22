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
	server.PubGet(server.OpenApiPath("/gallery/image/download"), func(c *gin.Context) {
		controller.DownloadImageThumbnailEndpoint(c, common.NewExecContext(c.Request.Context(), nil))
	})
	reportPath(ec, gclient.CreatePathReq{Url: server.OpenApiPath("/gallery/image/download"), Type: gclient.PT_PUBLIC,
		Desc: "Download gallery image by token", Method: "GET"})

	// authenticated routes
	server.Get(server.OpenApiPath("/gallery/brief/owned"), controller.ListOwnedGalleryBriefsEndpoint)
	reportPath(ec, gclient.CreatePathReq{Url: server.OpenApiPath("/gallery/brief/owned"), Type: gclient.PT_PROTECTED,
		Desc: "List owned gallery brief info", Method: "GET"})

	server.PostJ(server.OpenApiPath("/gallery/new"), controller.CreateGalleryEndpoint)
	reportPath(ec, gclient.CreatePathReq{Url: server.OpenApiPath("/gallery/new"), Type: gclient.PT_PROTECTED, Desc: "Create new gallery",
		Method: "POST"})

	server.PostJ(server.OpenApiPath("/gallery/update"), controller.UpdateGalleryEndpoint)
	reportPath(ec, gclient.CreatePathReq{Url: server.OpenApiPath("/gallery/update"), Type: gclient.PT_PROTECTED, Desc: "Update gallery",
		Method: "POST"})

	server.PostJ(server.OpenApiPath("/gallery/delete"), controller.DeleteGalleryEndpoint)
	reportPath(ec, gclient.CreatePathReq{Url: server.OpenApiPath("/gallery/delete"), Type: gclient.PT_PROTECTED, Desc: "Delete gallery",
		Method: "POST"})

	server.PostJ(server.OpenApiPath("/gallery/list"), controller.ListGalleriesEndpoint)
	reportPath(ec, gclient.CreatePathReq{Url: server.OpenApiPath("/gallery/list"), Type: gclient.PT_PROTECTED, Desc: "List galleries",
		Method: "POST"})

	server.PostJ(server.OpenApiPath("/gallery/access/grant"), controller.GrantGalleryAccessEndpoint)
	reportPath(ec, gclient.CreatePathReq{Url: server.OpenApiPath("/gallery/access/grant"), Type: gclient.PT_PROTECTED, Desc: "List granted access to the galleries",
		Method: "POST"})

	server.PostJ(server.OpenApiPath("/gallery/images"), controller.ListImagesEndpoint)
	reportPath(ec, gclient.CreatePathReq{Url: server.OpenApiPath("/gallery/images"), Type: gclient.PT_PROTECTED, Desc: "List images of gallery",
		Method: "POST"})

	server.PostJ(server.OpenApiPath("/gallery/image/transfer"), controller.TransferGalleryImageEndpoint)
	reportPath(ec, gclient.CreatePathReq{Url: server.OpenApiPath("/gallery/image/transfer"), Type: gclient.PT_PROTECTED,
		Desc:   "Host selected images on gallery",
		Method: "POST"})

	// bootstrap server
	server.DefaultBootstrapServer(os.Args)
}

func reportPath(ec common.ExecContext, r gclient.CreatePathReq) {
	server.OnServerBootstrapped(func() {
		r.Url = "/fantahsea" + r.Url
		r.Group = "fantahsea"
		r.ResCode = MNG_FILE_CODE
		if e := gclient.AddPath(ec.Ctx, r); e != nil {
			log.Fatalf("gclient.AddPath, req: %+v, %v", r, e)
		}
	})
}
