package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Chateaubriand-g/bili/auth_service/dao"
	"github.com/Chateaubriand-g/bili/auth_service/util"

	//"github.com/Chateaubriand-g/bili/auth_service/model"
	"github.com/Chateaubriand-g/bili/common/model"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct{ dao dao.UserDAO }

func NewAuthController(dao dao.UserDAO) *AuthController { return &AuthController{dao: dao} }

// Register 用户注册接口
// @Summary 用户注册
// @Description 接收用户名、密码、邮箱，完成用户注册
// @Tags auth
// @Accept application/json
// @Produce application/json
// @Param data body RegisterRequest true "注册请求参数"
// @Success 200 {object} model.APIResponse{data=RegisterResponse} "注册成功"
// @Success 400 {object} model.APIResponse "请求参数错误"
// @success 500 {object} model.APIResponse "服务器内部错误"
// @Router /accounts [post]
func (ctl *AuthController) Register(c *gin.Context) {
	log.Println("receiver request", c.Request.URL.Path)
	var in struct{ Username, Password, Email string }
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pw, _ := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	user := model.User{
		UserName: in.Username,
		PassWord: string(pw),
		Email:    in.Email,
	}
	if err := ctl.dao.Create(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"code": 200, "message": "注册成功"}})
}

// Login 用户登录接口
// @Summary 用户登录
// @Description 用户名密码验证通过后，生成JWT令牌并放回用户基础信息
// @Tags auth
// @Accept application/json
// @Produce application/json
// @Param data body RegisterRequest true "登录请求参数"
// @Success 200 {object} model.APIResponse{data=RegisterResponse} "登录成功(含token)"
// @Success 400 {object} model.APIResponse "请求参数错误或密码错误"
// @success 500 {object} model.APIResponse "服务器内部错误"
// @Router /auth/login [post]
func (ctl *AuthController) Login(c *gin.Context) {
	var in struct{ Username, Password string }
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user *model.User
	user, err := ctl.dao.FindByUsername(in.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PassWord), []byte(in.Password)) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password error"})
		return
	}
	accesstoken, _ := util.GenerateAccessToken(uint64(user.ID), "abc")
	refreshtoken, _ := util.GenerateRefreshToken()

	c.SetCookie(
		"refresh_token",
		refreshtoken,
		24*3600,
		"/",
		"175.178.78.121",
		false,
		true,
	)

	c.JSON(http.StatusOK, model.SuccessResponse(gin.H{"token": accesstoken, "user": gin.H{
		"username": user.UserName,
		"nickname": user.NickName,
	}}))
}

// Login 用户登出接口
// @Summary 用户登出
// @Description 客户端登出
// @Tags auth
// @Produce application/json
// @Success 200 {object} model.APIResponse{data=RegisterResponse} "登出成功"
// @Router /auth/logout [post]
func (ctl *AuthController) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (ctl *AuthController) RenewToken(c *gin.Context) {
	uidStr := c.GetHeader("X-User-ID")
	if uidStr == "" {
		c.JSON(http.StatusUnauthorized, model.BadResponse(400, "unauthorized", nil))
		return
	}
	uid, _ := strconv.ParseUint(uidStr, 10, 64)

	accessToken, _ := util.GenerateAccessToken(uid, "abc")

	c.JSON(http.StatusOK, model.SuccessResponse(gin.H{"token": accessToken}))
}

// 以下为Swagger文档专用请求/响应结构体（仅用于文档生成，无需实际业务处理）
// RegisterRequest 注册请求参数结构体
type RegisterRequest struct {
	Username string `json:"username" example:"test_user"` // 用户名（唯一）
	Password string `json:"password" example:"123456aB!"` // 密码（建议包含大小写字母+数字+符号）
	Email    string `json:"email" example:"test@xxx.com"` // 邮箱（用于后续验证/找回密码）
}

// RegisterResponse 注册响应数据结构体
type RegisterResponse struct {
	Code    int    `json:"code" example:"200"`     // 状态码
	Message string `json:"message" example:"注册成功"` // 提示信息
}

// LoginRequest 登录请求参数结构体
type LoginRequest struct {
	Username string `json:"username" example:"test_user"` // 用户名
	Password string `json:"password" example:"123456aB!"` // 密码
}

// LoginResponse 登录响应数据结构体
type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."` // JWT令牌（后续接口请求需携带）
	User  struct {
		Username string `json:"username" example:"test_user"` // 用户名
		Nickname string `json:"nickname" example:"测试用户"`      // 昵称（未设置则返回用户名）
	} `json:"user"` // 用户基础信息
}

// LogoutResponse 登出响应数据结构体
type LogoutResponse struct {
	Ok bool `json:"ok" example:"true"` // 登出状态（true表示成功）
}
