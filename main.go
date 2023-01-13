package main

import (
	"os"

	"github.com/curtisnewbie/fantahsea/data"
	"github.com/curtisnewbie/fantahsea/web/controller"
	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/server"
	"github.com/gin-gonic/gin"
)

func main() {
	common.DefaultReadConfig(os.Args)

	// register jobs
	common.ScheduleCron("0 0/10 * * * *", data.CleanUpDeletedGallery)
	common.GetScheduler().StartAsync()

	// routes
	server.PubGet(server.OpenApiPath("/test"), func(ctx *gin.Context) {
		common.WithTrace(ctx.Request.Context()).Infof("%v", ctx.Request.Header)
	})
	server.PubGet(server.OpenApiPath("/gallery/image/download"), controller.DownloadImageEndpoint)
	server.Get(server.OpenApiPath("/gallery/brief/owned"), server.BuildAuthRouteHandler(controller.ListOwnedGalleryBriefsEndpoint))
	server.Post(server.OpenApiPath("/gallery/new"), server.BuildAuthRouteHandler(controller.CreateGalleryEndpoint))
	server.Post(server.OpenApiPath("/gallery/update"), server.BuildAuthRouteHandler(controller.UpdateGalleryEndpoint))
	server.Post(server.OpenApiPath("/gallery/delete"), server.BuildAuthRouteHandler(controller.DeleteGalleryEndpoint))
	server.Post(server.OpenApiPath("/gallery/list"), server.BuildAuthRouteHandler(controller.ListGalleriesEndpoint))
	server.Post(server.OpenApiPath("/gallery/access/grant"), server.BuildAuthRouteHandler(controller.GrantGalleryAccessEndpoint))
	server.Post(server.OpenApiPath("/gallery/images"), server.BuildAuthRouteHandler(controller.ListImagesEndpoint))
	server.Post(server.OpenApiPath("/gallery/image/transfer"), server.BuildAuthRouteHandler(controller.TransferGalleryImageEndpoint))
	server.Post(server.OpenApiPath("/gallery/image/dir/transfer"), server.BuildAuthRouteHandler(controller.TransferGalleryImageInDir))

	// server
	server.BootstrapServer()
}
