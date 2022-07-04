package controller

import (
	"fantahsea/web/dto"

	"net/http"

	"github.com/gin-gonic/gin"
)

func GetDummy(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, dto.Dummy{ ID: "1", Title: "Good Dummy" })
}