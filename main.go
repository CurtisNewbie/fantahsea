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
	server.AddShutdownHook(func() { common.GetScheduler().Stop() })

	// public routes
	server.PubGet(server.OpenApiPath("/gallery/image/download"), func(c *gin.Context) {
		controller.DownloadImageThumbnailEndpoint(c, server.NewExecContext(c.Request.Context(), nil))
	})

	// authenticated routes
	server.Get(server.OpenApiPath("/gallery/brief/owned"), controller.ListOwnedGalleryBriefsEndpoint)
	server.Post(server.OpenApiPath("/gallery/new"), controller.CreateGalleryEndpoint)
	server.Post(server.OpenApiPath("/gallery/update"), controller.UpdateGalleryEndpoint)
	server.Post(server.OpenApiPath("/gallery/delete"), controller.DeleteGalleryEndpoint)
	server.Post(server.OpenApiPath("/gallery/list"), controller.ListGalleriesEndpoint)
	server.Post(server.OpenApiPath("/gallery/access/grant"), controller.GrantGalleryAccessEndpoint)
	server.Post(server.OpenApiPath("/gallery/images"), controller.ListImagesEndpoint)
	server.Post(server.OpenApiPath("/gallery/image/transfer"), controller.TransferGalleryImageEndpoint)
	server.Post(server.OpenApiPath("/gallery/image/dir/transfer"), controller.TransferGalleryImageInDir)

	// server
	server.BootstrapServer()
}
