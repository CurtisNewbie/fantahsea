package controller

import (
	"fmt"

	"github.com/curtisnewbie/fantahsea/config"
	"github.com/curtisnewbie/fantahsea/util"
	"github.com/curtisnewbie/fantahsea/weberr"

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

	router.Use(gin.CustomRecovery(func(c *gin.Context, e interface{}) {
		if err, ok := e.(error); ok {
			util.DispatchErrJson(c, err)
			return
		}
		util.DispatchErrJson(c, weberr.NewWebErr("Unknown error, please try again later"))
	}))

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
		return "open/api/" + relPath
	}

	return "remote" + relPath
}
