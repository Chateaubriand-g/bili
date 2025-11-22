package internal

import (
	"github.com/Chateaubriand-g/bili/common/mq"
	"github.com/gin-gonic/gin"
)

type AnalyticsController struct {
	dao AnalyticsDAO
}

func NewAnalyticsController(dao AnalyticsDAO, producer *mq.RocketMQProducer) *AnalyticsController {
	return &AnalyticsController{dao: dao}
}

func (ctl *AnalyticsController) GetVideoStats(c *gin.Context) {

}
