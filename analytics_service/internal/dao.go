package internal

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AnalyticsDAO interface {
	GetVideoDailyStats(videoID string) (map[string]int, error)
}

type analyticsDAO struct {
	DB  *gorm.DB
	RDS *redis.Client
}

func NewAnalyticsDAO(db *gorm.DB, rds *redis.Client) AnalyticsDAO {
	return &analyticsDAO{DB: db, RDS: rds}
}

func (dao *analyticsDAO) GetVideoDailyStats(videoID string) (map[string]int, error) {
	return nil, nil
}
