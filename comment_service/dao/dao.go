package dao

import (
	"github.com/Chateaubriand-g/bili/common/model"
	"github.com/redis/go-redis"
	"gorm.io/gorm"
)

type CommentDAO interface {
	AddComment(*model.Comment)
	IncrLikeNum()
	DecrLikeNum()
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

func (dao *commentDAO) AddComment(*model.Comment)
