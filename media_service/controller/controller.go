package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Chateaubriand-g/bili/common/alioss"
	"github.com/Chateaubriand-g/bili/common/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MediaController struct {
	DB   *gorm.DB
	OSSC *alioss.OssClient
	STSC *alioss.STSClient
}

func NewMediaController(db *gorm.DB, oss *alioss.OssClient, sts *alioss.STSClient) *MediaController {
	return &MediaController{DB: db, OSSC: oss, STSC: sts}
}

func (ctl *MediaController) GetSTS(c *gin.Context) {
	uidstr := c.GetHeader("X-User-ID")
	if uidstr == "" {
		c.JSON(http.StatusUnauthorized, model.BadResponse(400, "unauthorized", nil))
		return
	}

	sessionName := uidstr + fmt.Sprintf(":%d", time.Now().Unix())
	response, err := ctl.STSC.GetSTS(sessionName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.BadResponse(500, "getsts error: "+err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(response))
}
