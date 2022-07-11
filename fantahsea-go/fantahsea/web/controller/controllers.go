package controller

import (
	"fantahsea/config"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

const (
	ServerName = "fantahsea"	
)

func BootstrapServer(serverConf *config.ServerConfig, isProd bool) error {

	if isProd {
		log.Info("Using prod profile, will run with ReleaseMode")
		gin.SetMode(gin.ReleaseMode)
	}

	// register routes
	router := gin.Default()
	RegisterGalleryRoutes(router)
	// todo register routes for galleryImages

	// start the server
	err := router.Run(fmt.Sprintf("%v:%v", serverConf.Host, serverConf.Port))
	if err != nil {
		log.Errorf("Failed to bootstrap gin engine (web server), %v", err)
		return err
	}

	log.Printf("Web server bootstrapped on port: %v\n", serverConf.Port)

	return nil
}
