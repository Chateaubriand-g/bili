package controller

import (
	"net/http"
	"strconv"

	"github.com/Chateaubriand-g/bili/common/model"
	"github.com/Chateaubriand-g/bili/common/mq"
	"github.com/Chateaubriand-g/bili/interaction_service/dao"
	"github.com/gin-gonic/gin"
)

type InteractionController struct {
	dao      dao.InteractionDAO
	producer *mq.RocketMQProducer
}

func NewInteractionController(d dao.InteractionDAO, p *mq.RocketMQProducer) *InteractionController {
	return &InteractionController{
		dao:      d,
		producer: p,
	}
}

func (ctl *InteractionController) Like(c *gin.Context) {
	uidStr := c.GetHeader("X-User-ID")
	if uidStr == "" {
		c.JSON(http.StatusUnauthorized, model.BadResponse(400, "unauthorized", nil))
		return
	}
	uid, _ := strconv.ParseUint(uidStr, 10, 64)

	var in model.VideoLikeReq
	if err := c.ShouldBindBodyWithJSON(&in); err != nil {
		c.JSON(http.StatusUnauthorized, model.BadResponse(400, "required videolikereq", nil))
		return
	}

	switch in.Action {
	case "like":
		if err := ctl.dao.IncrLike(uid, in.VideoID); err != nil {
			c.JSON(http.StatusInternalServerError, model.BadResponse(500, err.Error(), nil))
			return
		}
	case "unlike":
		if err := ctl.dao.DecrLike(uid, in.VideoID); err != nil {
			c.JSON(http.StatusInternalServerError, model.BadResponse(500, err.Error(), nil))
			return
		}
	default:
		c.JSON(http.StatusInternalServerError, model.BadResponse(500, "unsupport action", nil))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil))
}

func (ctl *InteractionController) AddToFolder(c *gin.Context) {
	uidStr := c.GetHeader("X-User-ID")
	if uidStr == "" {
		c.JSON(http.StatusUnauthorized, model.BadResponse(400, "unauthorized", nil))
		return
	}
	//uid := strconv.ParseUint(uidStr, 10, 64)

	folderIDStr := c.Param("folder_id")
	videoIDStr := c.Param("video_id")
	if folderIDStr == "" || videoIDStr == "" {
		c.JSON(http.StatusBadRequest, model.BadResponse(400, "required param folder_id and video_id", nil))
		return
	}

	folderID, _ := strconv.ParseUint(folderIDStr, 10, 64)
	videoID, _ := strconv.ParseUint(videoIDStr, 10, 64)

	if err := ctl.dao.AddToFolder(videoID, folderID); err != nil {
		c.JSON(http.StatusBadRequest, model.BadResponse(500, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil))
}

func (ctl *InteractionController) DeleteFromFolder(c *gin.Context) {
	uidStr := c.GetHeader("X-User-ID")
	if uidStr == "" {
		c.JSON(http.StatusUnauthorized, model.BadResponse(400, "unauthorized", nil))
		return
	}
	//uid := strconv.ParseUint(uidStr, 10, 64)

	folderIDStr := c.Param("folder_id")
	videoIDStr := c.Param("video_id")
	if folderIDStr == "" || videoIDStr == "" {
		c.JSON(http.StatusBadRequest, model.BadResponse(400, "required param folder_id and video_id", nil))
		return
	}

	folderID, _ := strconv.ParseUint(folderIDStr, 10, 64)
	videoID, _ := strconv.ParseUint(videoIDStr, 10, 64)

	if err := ctl.dao.DeleteFromFolder(videoID, folderID); err != nil {
		c.JSON(http.StatusBadRequest, model.BadResponse(500, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil))
}

func (ctl *InteractionController) CreateFolder(c *gin.Context) {

}
