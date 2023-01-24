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
	// register jobs
	common.ScheduleCron("0 0/10 * * * *", data.CleanUpDeletedGallery)

	// public routes
	server.PubGet(server.OpenApiPath("/gallery/image/download"), func(c *gin.Context) {
		controller.DownloadImageThumbnailEndpoint(c, common.NewExecContext(c.Request.Context(), nil))
	})

	// authenticated routes
	server.Get(server.OpenApiPath("/gallery/brief/owned"), controller.ListOwnedGalleryBriefsEndpoint)
	server.PostJ(server.OpenApiPath("/gallery/new"), controller.CreateGalleryEndpoint)
	server.PostJ(server.OpenApiPath("/gallery/update"), controller.UpdateGalleryEndpoint)
	server.PostJ(server.OpenApiPath("/gallery/delete"), controller.DeleteGalleryEndpoint)
	server.PostJ(server.OpenApiPath("/gallery/list"), controller.ListGalleriesEndpoint)
	server.PostJ(server.OpenApiPath("/gallery/access/grant"), controller.GrantGalleryAccessEndpoint)
	server.PostJ(server.OpenApiPath("/gallery/images"), controller.ListImagesEndpoint)
	server.PostJ(server.OpenApiPath("/gallery/image/transfer"), controller.TransferGalleryImageEndpoint)

	// bootstrap server
	server.DefaultBootstrapServer(os.Args)
}
