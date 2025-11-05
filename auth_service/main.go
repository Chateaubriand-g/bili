package main

import (
	"log"

	"github.com/Chateaubriand-g/bili/auth_service/config"
	"github.com/Chateaubriand-g/bili/auth_service/controller"
	"github.com/Chateaubriand-g/bili/auth_service/dao"
	"github.com/Chateaubriand-g/bili/auth_service/middleware"
	"github.com/Chateaubriand-g/bili/auth_service/util"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("init config failed: %v", err)
	}

	deregiter, err := middleware.RegisterServiceToConsul(cfg)
	if err != nil {
		log.Fatalf("register service failed: %v", err)
	}
	defer deregiter()

	db, err := util.CreateDB(cfg)
	if err != nil {
		log.Fatalf("init databse failed: %v", err)
	}

	userDAO := dao.NewUserDAO(db)
	authCTL := controller.NewAuthController(userDAO)

	r := gin.Default()

	api := r.Group("/api")
	{
		api.POST("/user/account/register", authCTL.Register)
		api.POST("/user/account/login", authCTL.Login)
		api.POST("/user/account/logout", authCTL.Logout)
	}

	r.Run(":8081")
}
