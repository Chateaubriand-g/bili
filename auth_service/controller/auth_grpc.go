package controller

import (
	"context"

	"github.com/Chateaubriand-g/bili/auth_service/dao"
	"github.com/Chateaubriand-g/bili/pkg/authpb"
)

type GRPCServer struct {
	authpb.UnimplementedAuthServiceServer
	dao dao.UserDAO
}

func NewGRPCServer(dao dao.UserDAO) *GRPCServer {
	return &GRPCServer{dao: dao}
}

func (s *GRPCServer) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {

}
