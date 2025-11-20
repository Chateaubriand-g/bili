package dao

import (
	"gorm.io/gorm"
)

type InteractionDAO interface {
	CreateFolder(folderName string) error
	AddToFolder(videoID uint64) error
	DeleteFromFolder(videoID uint64) error
	IncrLike(userID, videID uint64) error
	DecrLike(userID, videoID uint64) error
}

type interactionDAO struct {
	DB *gorm.DB
}

func NewInteractionDAO(db *gorm.DB) InteractionDAO {
	return &interactionDAO{DB: db}
}

func (dao *interactionDAO) CreateFolder(folderName string) error {
	var newFolder 
}

func (dao *interactionDAO) AddToFolder(videoID uint64) error
func (dao *interactionDAO) DeleteFromFolder(videoID uint64) error
func (dao *interactionDAO) IncrLike(userID, videID uint64) error
func (dao *interactionDAO) DecrLike(userID, videoID uint64) error
