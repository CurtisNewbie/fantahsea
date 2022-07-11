package controller

import (
	"fantahsea/data"
	"fantahsea/util"
	"fantahsea/web/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterGalleryRoutes(router *gin.Engine) {
	router.PUT("/gallery/new", CreateGallery)
	router.POST("/gallery/update", UpdateGallery)
	router.POST("/gallery/delete", DeleteGallery)
}

/* CreateGallery web endpoint */
func CreateGallery(c *gin.Context) {
	user, err := util.ExtractUser(c)
	if err != nil {
		c.JSON(http.StatusOK, dto.WrapResp(nil, err))
	}

	var cmd *data.CreateGalleryCmd
	if err := c.ShouldBindJSON(cmd); err != nil {
		c.JSON(http.StatusOK, dto.ErrorResp("Illegal Arguments"))
		return
	}

	if _, err := data.CreateGallery(cmd, user); err != nil {
		c.JSON(http.StatusOK, dto.WrapResp(nil, err))
		return
	}

	c.JSON(http.StatusOK, dto.OkResp())
}

/* Update Gallery web endpoint */
func UpdateGallery(c *gin.Context) {
	user, err := util.ExtractUser(c)
	if err != nil {
		c.JSON(http.StatusOK, dto.WrapResp(nil, err))
	}

	var cmd *data.UpdateGalleryCmd
	if err := c.ShouldBindJSON(cmd); err != nil {
		c.JSON(http.StatusOK, dto.ErrorResp("Illegal Arguments"))
		return
	}

	if err := data.UpdateGallery(cmd, user); err != nil {
		c.JSON(http.StatusOK, dto.WrapResp(nil, err))
		return
	}

	c.JSON(http.StatusOK, dto.OkResp())
}

/* Delete Gallery web endpoint */
func DeleteGallery(c *gin.Context) {
	// todo
}
