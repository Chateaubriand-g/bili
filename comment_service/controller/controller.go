package controller

import (
	"context"
	"net/http"
	"strconv"

	"github.com/Chateaubriand-g/bili/comment_service/dao"
	"github.com/Chateaubriand-g/bili/common/middleware"
	"github.com/Chateaubriand-g/bili/common/model"
	"github.com/Chateaubriand-g/bili/common/mq"

	"github.com/gin-gonic/gin"
)

type CommentController struct {
	dao      dao.CommentDAO
	producer *mq.RocketMQProducer
}

func NewCommentController(dao dao.CommentDAO, p *mq.RocketMQProducer) *CommentController {
	return &CommentController{dao: dao, producer: p}
}

func (ctl *CommentController) GetCommentList(c *gin.Context) {
	uidstr := c.GetHeader("X-User-ID")
	if uidstr == "" {
		c.JSON(http.StatusUnauthorized, model.BadResponse(400, "unauthorized", nil))
		return
	}
	//uid, _ := strconv.ParseUint(uidstr, 10, 64)

	var in model.CommentListReq
	if err := c.ShouldBindBodyWithJSON(&in); err != nil {
		c.JSON(http.StatusUnauthorized, model.BadResponse(400, "required commentlistreq", nil))
		return
	}

	topList, err := ctl.dao.TopCommentList(CommentListReq.VideoID, CommentListReq.Page, CommentListReq.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.BadResponse(500, err.Error(), nil))
		return
	}

	secList, err := ctl.dao.SecCommentList(topList, CommentListReq.ReplySize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.BadResponse(500, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(gin.H{"comment": topList, "reply": secList}))
}

func (ctl *CommentController) AddComment(c *gin.Context) {
	uidstr := c.GetHeader("X-User-ID")
	if uidstr == "" {
		c.JSON(http.StatusUnauthorized, model.BadResponse(400, "unauthorized", nil))
		return
	}
	uid, _ := strconv.ParseUint(uidstr, 10, 64)

	var in model.CommentReq
	if err := c.ShouldBindBodyWithJSON(&in); err != nil {
		c.JSON(http.StatusUnauthorized, model.BadResponse(400, "struct commentreq required", nil))
		return
	}

	newComment := model.Comment{
		VideoID:  in.VideoID,
		UserID:   uid,
		Content:  in.Content,
		ParentID: in.ParentID,
	}
	commentID, err := ctl.dao.AddComment(&newComment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.BadResponse(500, "add comment failed", nil))
		return
	}

	toUserID, err := ctl.dao.FindUserIDByVideoID(in.VideoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.BadResponse(500, err.Error(), nil))
		return
	}

	v := middleware.MQMsg{
		UserID:     toUserID,
		Type:       0,
		FromUserID: uid,
		Payload: map[string]interface{}{
			"comment_id": commentID,
			"text":       in.Content,
		},
	}
	ctl.producer.SendEvent(context.TODO(), v)

	c.JSON(http.StatusOK, model.SuccessResponse(nil))
}

func (ctl *CommentController) ClickLike(c *gin.Context) {
	uidstr := c.GetHeader("X-User-ID")
	if uidstr == "" {
		c.JSON(http.StatusUnauthorized, model.BadResponse(400, "unauthorized", nil))
		return
	}
	uid, _ := strconv.ParseUint(uidstr, 10, 64)

	var in model.CommentLikeReq
	if err := c.ShouldBindBodyWithJSON(&in); err != nil {
		c.JSON(http.StatusUnauthorized, model.BadResponse(400, "unauthorized", nil))
		return
	}

	if in.Action == "like" {
		err := ctl.dao.IncrLikeNum(uid, in.CommentID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.BadResponse(500, "incrlikenum failed", nil))
			return
		}
	} else if in.Action == "unlike" {
		err := ctl.dao.DecrLikeNum(uid, in.CommentID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.BadResponse(500, "incrlikenum failed", nil))
			return
		}
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil))
}
