package dao

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const (
	minNotifyType = 0
	maxNotifyType = 5
)

type NotifyDAO interface {
	GetUnreadByType(userID uint64) (interface{}, error)
}

type notifyDAO struct {
	rds *redis.Client
	db  *gorm.DB
}

func NewNotifyDAO(db *gorm.DB, rds *redis.Client) NotifyDAO {
	return &notifyDAO{
		db:  db,
		rds: rds,
	}
}

func (dao *notifyDAO) GetUnreadByType(userID uint64) (interface{}, error) {
	typeKeys := make([]string, maxNotifyType-minNotifyType+1)
	for i := minNotifyType; i < maxNotifyType; i++ {
		typeKeys = append(typeKeys, fmt.Sprintf("notify:unread:type:%d:%d", userID, i))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	results, err := dao.rds.MGet(ctx, typeKeys...).Result()
	if err != nil {
		return nil, fmt.Errorf("rides mget error: %w", err)
	}

	counts := make([]uint64, maxNotifyType-minNotifyType+1)
	for i, _ := range typeKeys {
		val := results[i]
		if val == nil {
			//TODO
			counts[i] = 1
			continue
		}

		//results 返回的是interface{}类型
		valstr, ok := val.(string)
		if !ok {
			//TODO
			counts[i] = 1
			continue
		}

		typeCount, err := strconv.ParseUint(valstr, 10, 64)
		if err != nil {
			counts[i] = 1
			continue
		}

		counts[i] = typeCount
	}

	return counts, nil
}
