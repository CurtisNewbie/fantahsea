package controller

import (
	"fantahsea/data"
	"fantahsea/util"
	"fantahsea/web/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterGalleryRoutes(router *gin.Engine) {
	router.PUT(ServerName + "/gallery/new", CreateGallery)
	router.POST(ServerName + "/gallery/update", UpdateGallery)
	router.POST(ServerName + "/gallery/delete", DeleteGallery)
	router.POST(ServerName + "/gallery/list", ListGalleries)
}

/* ListGalleries web endpoint */
func ListGalleries(c *gin.Context) {
	user, err := util.ExtractUser(c)
	if err != nil {
		c.JSON(http.StatusOK, dto.WrapResp(nil, err))
		return
	}

	var page *dto.Paging
	if err := c.ShouldBindJSON(page); err != nil {
		c.JSON(http.StatusOK, dto.ErrorResp("Illegal Arguments"))
		return
	}

	galleries, err := data.ListGalleries(page, user)	
	if err != nil {
		c.JSON(http.StatusOK, dto.WrapResp(nil, err))
		return
	}

	c.JSON(http.StatusOK, dto.OkRespWData(galleries))
}

/* CreateGallery web endpoint */
func CreateGallery(c *gin.Context) {
	user, err := util.ExtractUser(c)
	if err != nil {
		c.JSON(http.StatusOK, dto.WrapResp(nil, err))
		return
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
		return
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
