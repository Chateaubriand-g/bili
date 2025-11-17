package main

import (
	"log"

	"github.com/Chateaubriand-g/bili/auth_service/config"
	"github.com/Chateaubriand-g/bili/auth_service/controller"
	"github.com/Chateaubriand-g/bili/auth_service/dao"
	"github.com/Chateaubriand-g/bili/auth_service/middleware"
	"github.com/Chateaubriand-g/bili/auth_service/util"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title bili_auth_service
// @version 1.0
// @description 路由分发，统计鉴权
// @termsOfService http://hostip/

// @BasePath /api/auth

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("init config failed: %v", err)
	}

	tracer, repoter, err := middleware.InitZipkin(cfg)
	if err != nil {
		log.Fatalf("initZipkin failed: %v", err)
	}
	defer middleware.CloseZipkin(repoter)

	deregiter, err := middleware.RegisterServiceToConsul(cfg)
	if err != nil {
		log.Fatalf("register service failed: %v", err)
	}
	defer deregiter()

	db, err := util.InitDatabase(cfg)
	if err != nil {
		log.Fatalf("init databse failed: %v", err)
	}

	userDAO := dao.NewUserDAO(db)
	authCTL := controller.NewAuthController(userDAO)

	r := gin.Default()

	if tracer != nil {
		r.Use(middleware.ZipkinMiddleware(tracer))
	}

	api := r.Group("/api")
	{
		api.POST("/auth/account/register", authCTL.Register)
		api.POST("/auth/account/login", authCTL.Login)
		api.POST("/auth/account/logout", authCTL.Logout)
	}

	r.GET("/swagger/*proxy", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8081")
}
