package util

import (
	"fmt"

	"github.com/Chateaubriand-g/bili/common/config"
	"github.com/Chateaubriand-g/bili/common/database"
	"gorm.io/gorm"
)

func InitDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/biliuser?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
	)
	return database.CreateDB(dsn)
}
