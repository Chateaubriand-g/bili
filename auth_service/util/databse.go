package util

import (
	"fmt"

	"github.com/Chateaubriand-g/bili/auth_service/config"
	"github.com/Chateaubriand-g/bili/common"

	"gorm.io/gorm"
)

func InitDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/biliuser?charset=utf8mb4&parseTime=True&loc=Local",
		config.DB.User,
		config.DB.Password,
		config.DB.Host,
		config.DB.Port,
	)
	return common.CreateDB()
}
