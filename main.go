package main

import (
	"os"

	"github.com/curtisnewbie/fantahsea/web/controller"
	"github.com/curtisnewbie/gocommon/web/server"
	"github.com/gin-gonic/gin"

	"github.com/curtisnewbie/gocommon/config"

	log "github.com/sirupsen/logrus"
)

func main() {

	profile := config.ParseProfile(os.Args[1:])
	log.Printf("Using profile: %v", profile)

	configFile := config.ParseConfigFilePath(os.Args[1:], profile)
	log.Printf("Looking for config file: %v", configFile)

	conf, err := config.ParseJsonConfig(configFile)
	if err != nil {
		panic(err)
	}
	config.SetGlobalConfig(conf)

	if err := config.InitDBFromConfig(&conf.DBConf); err != nil {
		panic(err)
	}

	isProd := profile == "prod"
	err = server.BootstrapServer(&conf.ServerConf, isProd, func(router *gin.Engine) {
		controller.RegisterGalleryRoutes(router)
		controller.RegisterGalleryImageRoutes(router)
	})
	if err != nil {
		panic(err)
	}

}
