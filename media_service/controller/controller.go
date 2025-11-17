package controller

import (
	sts "github.com/alibabacloud-go/sts-20150401/v2/client"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MediaController struct {
	DB  *gorm.DB
	OSS *oss.Client
	STS *sts.Client
}

func NewMediaController(db *gorm.DB, oss *oss.Client, sts *sts.Client) *MediaController {
	return &MediaController{DB: db, OSS: oss, STS: sts}
}

func (ctl *MediaController) GetSTS(c *gin.Context) {

}
