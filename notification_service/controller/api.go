package controller

import (
	"net/http"
	"strconv"

	"github.com/Chateaubriand-g/bili/common/model"
	"github.com/Chateaubriand-g/bili/notification_service/dao"
	"github.com/gin-gonic/gin"
)

type NotifyController struct {
	dao dao.NotifyDAO
}

func NewNotifyController(dao *dao.NotifyDAO) *NotifyController {
	return &NotifyController{
		dao: *dao,
	}
}

func (ctl *NotifyController) GetUnreadByType(c *gin.Context) {
	uidstr := c.GetHeader("X-User-ID")
	if uidstr == "" {
		c.JSON(http.StatusUnauthorized, model.BadResponse(400, "unauthorized", nil))
		return
	}

	uid, _ := strconv.ParseUint(uidstr, 10, 64)
	countTypes, err := ctl.dao.GetUnreadByType(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.BadResponse(500, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse(gin.H{
		"reply":    countTypes.([]int)[0],
		"at":       countTypes.([]int)[1],
		"love":     countTypes.([]int)[2],
		"system":   countTypes.([]int)[3],
		"whisper":  countTypes.([]int)[4],
		"dynamicL": countTypes.([]int)[5],
	}))
}
