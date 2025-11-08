package controller

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Chateaubriand-g/bili/common/model"
	"github.com/Chateaubriand-g/bili/user_service/dao"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	dao dao.UserDAO
	oss *oss.Client
}

func NewUserDAO(dao dao.UserDAO, oss *oss.Client) *UserController {
	return &UserController{dao: dao, oss: oss}
}

func (ctl *UserController) GetPersonalInfo(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, model.BadResponse(401, "unauthorized", nil))
		return
	}

	id, _ := strconv.ParseInt(userID, 10, 64)

	user, err := ctl.dao.FindByUserID(int(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.BadResponse(500, err.Error(), nil))
		return
	}

	var userDTO = model.UserDTO{
		ID:          user.ID,
		UserName:    user.UserName,
		NickName:    user.NickName,
		Gender:      user.Gender,
		Description: user.Description,
		Email:       user.Email,
		Avatar:      user.Avatar,
	}

	c.JSON(http.StatusOK, model.SuccessResponse(userDTO))
}

func (ctl *UserController) UpdateInfo(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, model.BadResponse(401, "unauthorized", nil))
		return
	}

	id, _ := strconv.ParseInt(userID, 10, 64)
	var in model.UserInfoUpdate
	if err := c.ShouldBindBodyWithJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, model.BadResponse(400, "unauthorized", nil))
		return
	}

	err := ctl.dao.Update(int(id), &in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.BadResponse(500, "sql update failed", nil))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponseWithMsg("ok", nil))
}

func (ctl *UserController) UpdateAvatar(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, model.BadResponse(401, "unauthorized", nil))
		return
	}

	id, _ := strconv.ParseInt(userID, 10, 64)
	avatarfile, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, model.BadResponse(400, "file required", nil))
	}

	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
	}

	fileType := avatarfile.Header.Get("Content-Type")
	if !allowedTypes[fileType] {
		c.JSON(http.StatusBadRequest, model.BadResponse(400, "unsupported file type (only jpg/png)", nil))
		return
	}

	f, _ := avatarfile.Open()
	defer f.Close()

	timestamp := time.Now().Format("20060102150405")
	tailfix := ".img"
	if fileType == "image/png" {
		tailfix = ".png"
	}
	key := fmt.Sprintf("%s_avatar_%s%s", userID, timestamp, tailfix)
	ossrequest := &oss.PutObjectRequest{
		Bucket:      oss.Ptr("bili-gz-a1"),
		Key:         oss.Ptr(key),
		Body:        f,
		ContentType: oss.Ptr(fileType),
	}

	_, err = ctl.oss.PutObject(context.TODO(), ossrequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.BadResponse(500, "failed to upload to OSS:"+err.Error(), nil))
		return
	}

	avatarURL := fmt.Sprintf("https://bili-user.oss-cn-beijing.aliyuncs.com/%s", key)
	var temp = model.UserAvatatUpdate{
		Avatar: avatarURL,
	}
	err = ctl.dao.Update(int(id), temp)
	if err != nil {
		// 这里可以考虑删除已上传的OSS文件（避免垃圾文件）
		_, _ = ctl.oss.DeleteObject(context.TODO(), &oss.DeleteObjectRequest{
			Bucket: oss.Ptr("bili-user"),
			Key:    oss.Ptr(key),
		})
		c.JSON(http.StatusInternalServerError, model.BadResponse(500, "failed to update user avatar: "+err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(temp))
}
