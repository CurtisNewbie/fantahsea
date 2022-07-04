package main

import (
	"fantahsea/config"
	"fantahsea/web/controller"
	"log"

	"github.com/gin-gonic/gin"
)


func main() {

	if err := config.InitDb(); err != nil {
		log.Fatalf("Failed to Open DB Handle, %v, exiting", err)
	}

	log.Println("Initializing Gin...")
	router := gin.Default()
	router.GET("/dummy", controller.GetDummy)
	router.Run("localhost:8080")
}