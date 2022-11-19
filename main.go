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
	server.AddRoutesRegistar(func(router *gin.Engine) {
		controller.RegisterGalleryRoutes(router)
		controller.RegisterGalleryImageRoutes(router)
	})
	
	// server
	server.BootstrapServer()
}
