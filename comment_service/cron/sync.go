package cron

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type SyncWorker struct {
	rds      *redis.Client
	db       *gorm.DB
	interval time.Duration
	stopChan chan struct{}
}

func NewSyncWorker(rds *redis.Client, db *gorm.DB, interval time.Duration) *SyncWorker {
	return &SyncWorker{
		rds:      rds,
		db:       db,
		interval: interval,
		stopChan: make(chan struct{}),
	}
}

func (s *SyncWorker) Start() {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := s.syncCounts(); err != nil {
				log.Printf("sync failed: %v", err)
			}
		case <-s.stopChan:
			s.syncCounts()
			return
		}
	}
}

func (s *SyncWorker) Stop() {
	close(s.stopChan)
}

func (s *SyncWorker) syncCounts() error {
	lockKey := "sync:lock:comment"
	lock, err := s.rds.SetNX(context.TODO(), lockKey, "1", 10*time.Minute).Result()
	if err != nil {
		return err
	}

	if !lock {
		return nil
	}
	defer s.rds.Del(context.TODO(), lockKey)

	dityKey := "comment:dirty"
	ids, err := s.rds.SMembers(context.TODO(), dityKey).Result()
	if err != nil {
		return err
	}

	if len(ids) == 0 {
		return nil
	}

	pipe := s.rds.Pipeline()
	countCmds := make(map[string]*redis.StringCmd)

	for _, id := range ids {
		countCmds[id] = pipe.Get(context.TODO(), "comment:count:"+id)
	}

	if _, err := pipe.Exec(context.TODO()); err != nil && err != redis.Nil {
		return err
	}

	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	for _, id := range ids {
		cmd := countCmds[id]
		count, err := cmd.Int()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				count = 0
			} else {
				tx.Rollback()
				return err
			}
		}

		if err := tx.Exec("update video SET commnt_count = ? where id = ?", count, id).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	if _, err := s.rds.SRem(context.Background(), dityKey, ids).Result(); err != nil {
		return err
	}

	return nil
}
