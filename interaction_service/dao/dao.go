package dao

import (
	"context"
	"fmt"

	"github.com/Chateaubriand-g/bili/common/model"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type InteractionDAO interface {
	CreateFolder(folderName string) error
	AddToFolder(videoID, folderID uint64) error
	DeleteFromFolder(videoID, folderID uint64) error
	IncrLike(userID, videoID uint64) error
	DecrLike(userID, videoID uint64) error
}

type interactionDAO struct {
	DB  *gorm.DB
	RDS *redis.Client
}

func NewInteractionDAO(db *gorm.DB, rds *redis.Client) InteractionDAO {
	return &interactionDAO{DB: db, RDS: rds}
}

func (dao *interactionDAO) CreateFolder(userID uint64, folderName string) error {
	newFolder := model.Folder{
		UserID:     userID,
		FolderName: folderName,
	}

	if err := dao.DB.Create(&newFolder).Error; err != nil {
		return fmt.Errorf("create folder error: %w", err)
	}
	return nil
}

func (dao *interactionDAO) AddToFolder(videoID, folderID uint64) error {
	newItem := model.FolderItems{
		FloderID: folderID,
		VideoID:  videoID,
	}

	if err := dao.DB.Create(&newItem).Error; err != nil {
		return fmt.Errorf("add to folder error: %w", err)
	}
	return nil
}

func (dao *interactionDAO) DeleteFromFolder(videoID, folderID uint64) error {
	target := model.FolderItems{
		FolderID: folderID,
		VideoID:  videoID,
	}

	if err := dao.DB.Delete(&target).Error; err != nil {
		return fmt.Errorf("del from folder error: %w", err)
	}
	return nil
}

func (dao *interactionDAO) IncrLike(userID, videoID uint64) error {
	newVideoLike := model.VideoLike{
		UserID:  userID,
		VideoID: videoID,
	}

	if err := dao.DB.Create(&newVideoLike).Error; err != nil {
		return fmt.Errorf("incrlike error: %w", err)
	}

	key := fmt.Sprintf("video:count:like:%d", videoID)
	dao.RDS.Incr(context.TODO(), key)
	dao.RDS.SAdd(context.TODO(), "video::dirty", videoID)

	return nil
}

func (dao *interactionDAO) DecrLike(userID, videoID uint64) error {
	if err := dao.DB.Delete(&model.VideoLike{}, "user_id = ? and video_id = ?", userID, videoID).Error; err != nil {
		return fmt.Errorf("decrlike error: %w", err)
	}
	return nil
}
