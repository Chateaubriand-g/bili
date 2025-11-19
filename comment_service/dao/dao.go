package dao

import (
	"fmt"
	"strconv"

	"github.com/Chateaubriand-g/bili/common/model"
	"github.com/redis/go-redis"
	"gorm.io/gorm"
)

type CommentDAO interface {
	AddComment(req *model.CommentReq) error
	IsLike(userID, commentID uint64) (bool, error)
	IncrLikeNum(userID, commentID uint64) error
	DecrLikeNum(userID, commentID uint64) error
}

type commentDAO struct {
	DB  *gorm.DB
	RDS *redis.Client
}

func NewCommentDAO(db *grom.DB, rds *redis.Client) CommentDAO {
	return &commentDAO{
		DB:  db,
		RDS: rds,
	}
}

func (dao *commentDAO) AddComment(req *model.CommentReq) error {
	newComm := model.Comment{
		VideoID:  req.VideoID,
		Content:  req.Content,
		ParentID: req.ParentID,
	}

	if err := dao.DB.Create(&newComm).Error; err != nil {
		return fmt.Errorf("db create err: %w", err)
	}

	key := "comment:count:" + strconv.Itoa(req.VideoID)
	dao.RDS.Incr(key)
	dao.RDS.SAdd("comment:dirty", req.VideoID)

	return nil
}

func (dao *commentDAO) IsLike(userID, commentID uint64) (bool, error) {
	var count int
	err := dao.DB.Model(&model.CommenstLike{}).Count("user_id = ? and comment_id = ?", userID, commentID).Error
	if err != nil {
		return false, fmt.Errorf("count failed: %w", err)
	}
	if count == 1 {
		return true, nil
	}
	return false, nil
}

func (dao *commentDAO) IncrLikeNum(userID, commentID uint64) error {
	newCommentLike := model.CommentsLike{
		CommentID: commentID,
		UserID:    userID,
	}
	if err := dao.DB.Create(&newCommentLike).Error; err != nil {
		return fmt.Errorf("create comments_like error: %w", err)
	}
	return nil
}

func (dao *commentDAO) DecrLikeNum(userID, commentID uint64) error {
	if err := dao.DB.Model(&model.CommentsLike).
		Delete("user_id = ? and comment_id = ?", userID, commentID).Error; err != nil {
		return fmt.Errorf("delete comments_like error: %w", err)
	}
	return nil
}
