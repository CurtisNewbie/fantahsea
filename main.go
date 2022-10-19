package main

import (
	"github.com/curtisnewbie/fantahsea/data"
	"github.com/curtisnewbie/fantahsea/web/controller"
	"github.com/curtisnewbie/gocommon/consul"
	"github.com/curtisnewbie/gocommon/util"
	"github.com/curtisnewbie/gocommon/web/server"
	"github.com/gin-gonic/gin"

	"github.com/curtisnewbie/gocommon/config"
)

func main() {

	_, conf := config.DefaultParseProfConf()

	// register jobs
	s := util.ScheduleCron("0 0/10 * * * *", data.CleanUpDeletedGallery)
	s.StartAsync()

	server.BootstrapServer(conf, func(router *gin.Engine) {
		consul.RegisterDefaultHealthCheck(router)
		controller.RegisterGalleryRoutes(router)
		controller.RegisterGalleryImageRoutes(router)
	})
}
