package main

import (
	"os"
	"strings"

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
	server.AddRoutesRegistar(func(router *gin.Engine) {
		controller.RegisterGalleryRoutes(router)
		controller.RegisterGalleryImageRoutes(router)
	})

	// whitelist for authorization
	server.AddRouteAuthWhitelist(func(url string) bool {
		return strings.HasPrefix(url, server.ResolvePath("/gallery/image/download", true))
	})

	// server
	server.BootstrapServer()
}
