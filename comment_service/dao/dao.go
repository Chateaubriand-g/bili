package dao

import (
	"context"
	"fmt"

	"github.com/Chateaubriand-g/bili/common/model"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type CommentDAO interface {
	TopCommentList(videoId, page, pageSize uint64) ([]*model.Comment, error)
	SecCommentList(toplist []*model.Comment, SecSize uint64) ([]*model.Comment, error)
	AddComment(req *model.Comment) (uint64, error)
	IsLike(userID, commentID uint64) (bool, error)
	IncrLikeNum(userID, commentID uint64) error
	DecrLikeNum(userID, commentID uint64) error
	FindUserIDByVideoID(videoID uint64) (uint64, error)
}

type commentDAO struct {
	DB  *gorm.DB
	RDS *redis.Client
}

func NewCommentDAO(db *gorm.DB, rds *redis.Client) CommentDAO {
	return &commentDAO{
		DB:  db,
		RDS: rds,
	}
}

func (dao *commentDAO) TopCommentList(videoId, page, pageSize uint64) ([]*model.Comment, error) {
	var ret []*model.Comment
	offset := (page - 1) * pageSize

	if err := dao.DB.Model(&model.Comment{}).
		Where("video_id = ? and parent_id = 0", videoId).Order("id DESC").
		Offset(int(offset)).Limit(int(pageSize)).Find(&ret).Error; err != nil {
		return nil, fmt.Errorf("top commentlist err: %w", err)
	}

	return ret, nil
}

func (dao *commentDAO) SecCommentList(toplist []*model.Comment, SecSize uint64) ([]*model.Comment, error) {
	ret := make([]*model.Comment, len(toplist))
	for _, value := range toplist {
		var temp *model.Comment
		if err := dao.DB.Where("parent_id = ?", value.ID).Limit(int(SecSize)).Find(&temp).Error; err != nil {
			ret = append(ret, &model.Comment{})
			continue
		}
		ret = append(ret, temp)
	}
	return ret, nil
}

func (dao *commentDAO) AddComment(req *model.Comment) (uint64, error) {
	if err := dao.DB.Create(&req).Error; err != nil {
		return 0, fmt.Errorf("db create err: %w", err)
	}

	key := fmt.Sprintf("comment:count:%d", req.VideoID)
	dao.RDS.Incr(context.TODO(), key)
	dao.RDS.SAdd(context.TODO(), "comment:dirty", req.VideoID)

	return req.ID, nil
}

func (dao *commentDAO) IsLike(userID, commentID uint64) (bool, error) {
	var count int64
	err := dao.DB.Model(&model.CommentsLike{}).
		Where("user_id = ? and comment_id = ?", userID, commentID).
		Count(&count).Error
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

	key := fmt.Sprintf("comment:likenums:%d", commentID)
	dao.RDS.Incr(context.TODO(), key)
	return nil
}

func (dao *commentDAO) DecrLikeNum(userID, commentID uint64) error {
	if err := dao.DB.Model(&model.CommentsLike{}).
		Delete("user_id = ? and comment_id = ?", userID, commentID).Error; err != nil {
		return fmt.Errorf("delete comments_like error: %w", err)
	}

	key := fmt.Sprintf("comment:likenums:%d", commentID)
	dao.RDS.Decr(context.TODO(), key)
	return nil
}

func (dao *commentDAO) FindUserIDByVideoID(videoId uint64) (uint64, error) {
	var userID uint64
	if err := dao.DB.Model(&model.Video{}).Where("video_id = ?", videoId).
		First(&userID).Error; err != nil {
		return 0, fmt.Errorf("find userid by videoid err: %w", err)
	}
	return userID, nil
}
