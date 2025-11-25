package controller

import (
	"context"
	"errors"
	"strconv"

	"github.com/Chateaubriand-g/bili/auth_service/dao"
	"github.com/Chateaubriand-g/bili/auth_service/util"
	"github.com/Chateaubriand-g/bili/common/middleware"
	"github.com/Chateaubriand-g/bili/common/model"
	authpb "github.com/Chateaubriand-g/bili/pkg/pb/auth"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type GRPCServer struct {
	authpb.UnimplementedAuthServiceServer
	dao dao.UserDAO
}

func NewGRPCServer(dao dao.UserDAO) *GRPCServer {
	return &GRPCServer{dao: dao}
}

func (s *GRPCServer) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	if req.GetEmail() == "" || req.GetPassword() == "" || req.GetUsername() == "" {
		return nil, status.Error(codes.InvalidArgument, "username,password and email are required")
	}

	hashpw, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), 12)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "hash password failed: %v", err)
	}

	newUser := &model.User{
		UserName: req.GetUsername(),
		PassWord: string(hashpw),
		Email:    req.GetEmail(),
	}

	err = s.dao.Create(newUser)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create newuser err: %v", err)
	}

	return &authpb.RegisterResponse{
		Userid: strconv.FormatUint(uint64(newUser.ID), 10),
	}, nil
}

func (s *GRPCServer) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	if req.GetUsername() == "" || req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "username and password are required")
	}
	user, err := s.dao.FindByUsername(req.GetUsername())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.Unauthenticated, "invaild username or password")
		}
		return nil, status.Errorf(codes.Internal, "query user failed: %v", err)
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PassWord), []byte(req.GetPassword())) != nil {
		return nil, status.Error(codes.Unauthenticated, "invaild username or password")
	}
	accesstoken, err := util.GenerateAccessToken(uint64(user.ID), "abc")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "generate token err: %v", err)
	}

	refreshtoken, err := util.GenerateRefreshToken()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "generate rtoken err: %v", err)
	}
	s.dao.SaveRefreshToken(refreshtoken, uint64(user.ID))

	return &authpb.LoginResponse{
		AccessToken:  accesstoken,
		RefreshToken: refreshtoken,
		Username:     user.UserName,
		Nickname:     user.NickName,
	}, nil
}

func (s *GRPCServer) Logout(ctx context.Context, req *authpb.LogoutRequest) (*authpb.LogoutResponse, error) {
	if req.GetUserid() == "" || req.GetRefreshToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "userid and refreshtoken is required")
	}
	uid, _ := strconv.ParseUint(req.GetUserid(), 10, 64)
	_ := s.dao.DeleteRefreshToken(req.GetRefreshToken(), uid)
	return &authpb.LogoutResponse{}, nil
}

func (s *GRPCServer) RefreshToken(ctx context.Context, req *authpb.RefreshTokenRequest) (*authpb.RefreshTokenResponse, error) {
	if req.GetUserid() == "" || req.GetRefreshToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "userid and refreshtoken is required")
	}

	uid, err := strconv.ParseUint(req.GetUserid(), 10, 64)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "userid must be numeric")
	}

	if !s.dao.IsTokenVaild(req.GetRefreshToken(), uid) {
		return nil, status.Error(codes.Internal, "invaild refreshtoken")
	}

	accessToken, err := middleware.GenerateAccessToken(uid, "abc")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "generate access token failed: %v", err)
	}

	return &authpb.RefreshTokenResponse{
		AccessToken: accessToken,
	}, nil
}
