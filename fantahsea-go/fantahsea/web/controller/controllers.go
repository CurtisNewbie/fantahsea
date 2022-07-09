package controller

import (
	"fantahsea/config"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func BootstrapServer(serverConf *config.ServerConfig) error {

	// register routes
	router := gin.Default()
	router.GET("/dummy", GetDummy)

	// start the server
	err := router.Run(fmt.Sprintf("%v:%v", serverConf.Host, serverConf.Port))
	if err != nil {
		log.Printf("Failed to run gin router, %v", err)
		return err
	}

	log.Printf("Web server bootstrapped on port: %v\n", serverConf.Port)

	return nil
}
