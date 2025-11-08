package controller

import (
	"net/http"

	"github.com/Chateaubriand-g/bili/user_service/dao"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	dao dao.UserDAO
}

func NewUserDAO(dao dao.UserDAO) *UserController { return &UserController{dao: dao} }

func (ctl *UserController) GetPersonalInfo(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized)
	}
}
