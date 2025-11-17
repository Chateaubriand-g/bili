package main

import (
	"github.com/Chateaubriand-g/bili/common/config"
	"github.com/Chateaubriand-g/bili/common/middleware"
	"github.com/Chateaubriand-g/bili/common/oss"
	"github.com/Chateaubriand-g/bili/media_service/util"
)

func main() {
	cfg, err := config.LoadConfig()

	db, err := util.InitDatabase(cfg)

	rds, err := middleware.InitRedis(cfg)

	deregister, err := middleware.RegisterServiceToConsul(cfg)
	defer deregister()

	tracer, reporter, err := middleware.InitZipkin(cfg)
	defer middleware.CloseZipkin(reporter)

	oss := oss.InitOssClient(cfg)
}
