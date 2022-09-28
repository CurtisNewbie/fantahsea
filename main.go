package main

import (
	"github.com/curtisnewbie/fantahsea/data"
	"github.com/curtisnewbie/fantahsea/web/controller"
	"github.com/curtisnewbie/gocommon/util"
	"github.com/curtisnewbie/gocommon/web/server"
	"github.com/gin-gonic/gin"

	"github.com/curtisnewbie/gocommon/config"
)

func main() {

	profile, conf := config.DefaultParseProfConf()

	if err := config.InitDBFromConfig(&conf.DBConf); err != nil {
		panic(err)
	}
	config.InitRedisFromConfig(&conf.RedisConf)

	// register jobs
	s := util.ScheduleCron("*/3 * * * * *", data.CleanUpDeletedGallery)
	s.StartAsync()

	isProd := config.IsProd(profile)
	err := server.BootstrapServer(&conf.ServerConf, isProd, func(router *gin.Engine) {
		controller.RegisterGalleryRoutes(router)
		controller.RegisterGalleryImageRoutes(router)
	})
	if err != nil {
		panic(err)
	}

}
