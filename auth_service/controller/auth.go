package controller

import (
	"log"
	"net/http"

	"github.com/Chateaubriand-g/bili/auth_service/dao"
	"github.com/Chateaubriand-g/bili/auth_service/middleware"
	"github.com/Chateaubriand-g/bili/auth_service/model"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct{ dao dao.UserDAO }

func NewAuthController(dao dao.UserDAO) *AuthController { return &AuthController{dao: dao} }

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
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (ctl *AuthController) Login(c *gin.Context) {
	var in struct{ Username, Password string }
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user model.User
	if _, err := ctl.dao.FindByUsername(in.Username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PassWord), []byte(in.Password)) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password error"})
		return
	}
	token, _ := middleware.GenerateToken(uint64(user.ID), user.UserName)
	c.JSON(http.StatusOK, gin.H{"token": token, "user": gin.H{"id": user.ID, "username": user.UserName}})
}

func (ctl *AuthController) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
