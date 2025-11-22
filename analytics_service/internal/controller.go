package internal

import "github.com/gin-gonic/gin"

type AnalyticsController struct {
	dao AnalyticsDAO
}

func NewAnalyticsController(dao AnalyticsDAO) *AnalyticsController {
	return &AnalyticsController{dao: dao}
}

func (ctl *AnalyticsController) GetVideoStats(c *gin.Context) {

}
