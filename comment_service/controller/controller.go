package controller

import (
	"net/http"

	"github.com/Chateaubriand-g/bili/comment_service/dao"
	"github.com/Chateaubriand-g/bili/common/model"
	"github.com/gin-gonic/gin"
)

type CommentController struct {
	dao dao.CommentDAO
}

func NewCommentController(dao dao.CommentDAO) *CommentController {
	return &CommentController{dao: dao}
}

func (ctl *CommentController) AddComment(c *gin.Context) {
	uidstr := c.GetHeader("X-User-ID")
	if uidstr == "" {
		c.JSON(http.StatusUnauthorized, model.BadResponse(400, "unauthorized", nil))
		return
	}

	var in model.CommentReq
	if err := c.ShouldBindBodyWithJSON(&in); err != nil {
		c.JSON(http.StatusUnauthorized, model.BadResponse(400, "struct commentreq required", nil))
		return
	}

	err := ctl.dao.AddComment(&in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.BadResponse(500, "add comment failed", nil))
		return
	}
}
