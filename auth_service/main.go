package main

import (
	"fmt"
	"log"
	"net"

	"github.com/Chateaubriand-g/bili/auth_service/config"
	"github.com/Chateaubriand-g/bili/auth_service/controller"
	"github.com/Chateaubriand-g/bili/auth_service/dao"
	"github.com/Chateaubriand-g/bili/auth_service/util"
	"github.com/Chateaubriand-g/bili/common/middleware"
	authpb "github.com/Chateaubriand-g/bili/pkg/pb/auth"
	"google.golang.org/grpc"
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

	// tracer, repoter, err := middleware.InitZipkin(cfg)
	// if err != nil {
	// 	log.Fatalf("initZipkin failed: %v", err)
	// }
	// defer middleware.CloseZipkin(repoter)

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
	// authCTL := controller.NewAuthController(userDAO)

	// r := router.InitRouter(authCTL, tracer)
	// r.Run(":8081")
	grpcServer := controller.NewGRPCServer(userDAO)
	listenPort := cfg.Server.Port
	if listenPort == 0 {
		listenPort = 8081
	}
	addr := fmt.Sprintf(":%d", listenPort)

	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen on %s failed: %v", addr, err)
	}

	server := grpc.NewServer()
	authpb.RegisterAuthServiceServer(server, grpcServer)

	server.Serve(listen)
}
