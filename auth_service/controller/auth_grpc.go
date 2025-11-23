package controller

import (
	"context"
	"net/http"

	"github.com/Chateaubriand-g/bili/auth_service/dao"
	"github.com/Chateaubriand-g/bili/auth_service/middleware"
	"github.com/Chateaubriand-g/bili/common/model"
	"github.com/Chateaubriand-g/bili/pkg/authpb"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type GRPCServer struct {
	authpb.UnimplementedAuthServiceServer
	dao dao.UserDAO
}

func NewGRPCServer(dao dao.UserDAO) *GRPCServer {
	return &GRPCServer{dao: dao}
}

func (s *GRPCServer) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	var user *model.User
	user, err := s.dao.FindByUsername(req.username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PassWord), []byte(in.Password)) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password error"})
		return
	}
	accesstoken, _ := middleware.GenerateAccessToken(uint64(user.ID), "abc")
	refreshtoken, _ := middleware.GenerateRefreshToken(uint64(user.ID), "abc")

}
