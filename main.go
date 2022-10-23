package main

import (
	"os"

	"github.com/curtisnewbie/fantahsea/data"
	"github.com/curtisnewbie/fantahsea/web/controller"
	"github.com/curtisnewbie/gocommon/consul"
	"github.com/curtisnewbie/gocommon/util"
	"github.com/curtisnewbie/gocommon/web/server"
	"github.com/gin-gonic/gin"

	"github.com/curtisnewbie/gocommon/config"
)

func main() {
	_, conf := config.DefaultParseProfConf(os.Args)

	// register jobs
	util.ScheduleCron("0 0/10 * * * *", data.CleanUpDeletedGallery)
	util.ScheduleCron("0 0/1 * * * *", consul.PollServiceListInstances)
	util.GetScheduler().StartAsync()

	server.BootstrapServer(conf, func(router *gin.Engine) {
		controller.RegisterGalleryRoutes(router)
		controller.RegisterGalleryImageRoutes(router)
	})
}
