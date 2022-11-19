package main

import (
	"os"

	"github.com/curtisnewbie/fantahsea/data"
	"github.com/curtisnewbie/fantahsea/web/controller"
	"github.com/curtisnewbie/gocommon"
	"github.com/gin-gonic/gin"
)

func main() {
	gocommon.DefaultReadConfig(os.Args)

	// init consul client and the polling server list subscription
	gocommon.GetConsulClient()

	// register jobs
	gocommon.ScheduleCron("0 0/10 * * * *", data.CleanUpDeletedGallery)
	gocommon.GetScheduler().StartAsync()

	// routes
	gocommon.AddRoutesRegistar(func(router *gin.Engine) {
		controller.RegisterGalleryRoutes(router)
		controller.RegisterGalleryImageRoutes(router)
	})
	
	// server
	gocommon.BootstrapServer()
}
