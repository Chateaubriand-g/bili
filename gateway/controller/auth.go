package controller

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/Chateaubriand-g/bili/common/model"
	"github.com/Chateaubriand-g/bili/gateway/proxy"
	authpb "github.com/Chateaubriand-g/bili/pkg/pb/auth"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/status"
)

type AuthController struct {
	consul *api.Client

	mtx       sync.RWMutex
	rpcConn   *grpc.ClientConn
	rpcClient authpb.AuthServiceClient
}

func NewAuthController(cli *api.Client) *AuthController {
	return &AuthController{consul: cli}
}

func (ctl *AuthController) Register(c *gin.Context) {
	var in struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindBodyWithJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, model.BadResponse(400, err.Error(), nil))
		return
	}

	client, err := ctl.getClient()
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	resp, err := client.Register(ctx, &authpb.RegisterRequest{
		Username: in.Username,
		Password: in.Password,
		Email:    in.Email,
	})
	if err != nil {
		ctl.handleRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(gin.H{"user_id": resp.GetUserid()}))
}

func (ctl *AuthController) Login(c *gin.Context) {
	var in struct{ Username, Password string }
	if err := c.ShouldBindBodyWithJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, model.BadResponse(400, err.Error(), nil))
		return
	}

	client, err := ctl.getClient()
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	resp, err := client.Login(ctx, &authpb.LoginRequest{
		Username: in.Username,
		Password: in.Password,
	})
	if err != nil {
		ctl.handleRPCError(c, err)
		return
	}

	c.SetCookie(
		"refresh_token",
		resp.GetRefreshToken(),
		24*3600,
		"/",
		"175.178.78.121",
		false,
		true,
	)

	c.JSON(http.StatusOK, model.SuccessResponse(gin.H{"token": resp.GetAccessToken(), "user": gin.H{
		"username": resp.GetUsername(),
		"nickname": resp.GetNickname(),
	}}))
}

func (ctl *AuthController) Logout(c *gin.Context) {
	client, err := ctl.getClient()
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	_, err = client.Logout(ctx, &authpb.LogoutRequest{
		Userid: c.GetHeader("X-User-ID"),
	})
	if err != nil {
		ctl.handleRPCError(c, err)
		return
	}

}

func (ctl *AuthController) getClient() (authpb.AuthServiceClient, error) {
	ctl.mtx.RLock()
	if ctl.rpcClient != nil && ctl.rpcConn != nil && ctl.rpcConn.GetState() != connectivity.Shutdown {
		client := ctl.rpcClient
		ctl.mtx.RUnlock()
		return client, nil
	}
	ctl.mtx.RUnlock()

	ctl.mtx.Lock()
	defer ctl.mtx.Unlock()

	if ctl.rpcClient != nil && ctl.rpcConn != nil && ctl.rpcConn.GetState() != connectivity.Shutdown {
		return ctl.rpcClient, nil
	}

	conn, err := proxy.GetGRPCConn(ctl.consul, "auth-service")
	if err != nil {
		return nil, err
	}

	ctl.rpcConn = conn
	ctl.rpcClient = authpb.NewAuthServiceClient(conn)
	return ctl.rpcClient, nil
}

func (ctl *AuthController) handleRPCError(c *gin.Context, err error) {
	st, ok := status.FromError(err)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	switch st.Code() {
	case codes.InvalidArgument:
		c.JSON(http.StatusBadRequest, gin.H{"error": st.Message()})
	case codes.AlreadyExists:
		c.JSON(http.StatusConflict, gin.H{"error": st.Message()})
	case codes.Unauthenticated:
		c.JSON(http.StatusUnauthorized, gin.H{"error": st.Message()})
	case codes.NotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": st.Message()})
	case codes.Unavailable:
		ctl.resetClient()
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": st.Message()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
	}
}

func (ctl *AuthController) resetClient() {
	ctl.mtx.Lock()
	defer ctl.mtx.Unlock()
	ctl.rpcClient = nil
	ctl.rpcConn = nil
}
