package controller

import (
	"fantahsea/config"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

/* Bootstrap Server */
func BootstrapServer(serverConf *config.ServerConfig, isProd bool) error {

	if isProd {
		log.Info("Using prod profile, will run with ReleaseMode")
		gin.SetMode(gin.ReleaseMode)
	}

	// register routes
	router := gin.Default()
	RegisterGalleryRoutes(router)
	RegisterGalleryImageRoutes(router)

	// start the server
	err := router.Run(fmt.Sprintf("%v:%v", serverConf.Host, serverConf.Port))
	if err != nil {
		log.Errorf("Failed to bootstrap gin engine (web server), %v", err)
		return err
	}

	log.Printf("Web server bootstrapped on port: %v", serverConf.Port)

	return nil
}

// Resolve request path
func ResolvePath(relPath string, isOpenApi bool) string {
	if isOpenApi {
		return "open" + relPath
	}

	return "remote" + relPath
}
