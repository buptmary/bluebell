package controller

import (
	"bluebell/dao/mysql"
	"bluebell/logic"
	"bluebell/models"
	"bluebell/pkg/jwt"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

func SignUpHandler(c *gin.Context) {
	// 1. 获取请求参数
	var fo models.RegisterForm

	// 2. 校验数据有效性
	if err := c.ShouldBindJSON(&fo); err != nil {
		ResponseErrorWithMsg(c, CodeInvalidParams, err.Error())
		return
	}

	// 3. 注册用户
	if err := logic.SignUp(&fo); err != nil {
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseErrorWithMsg(c, CodeServerBusy, err.Error())
		return
	}
	ResponseSuccess(c, nil)
}

func LoginHandler(c *gin.Context) {
	// 1. 获取请求参数及参数校验
	var fo models.LoginForm
	if err := c.ShouldBindJSON(&fo); err != nil {
		ResponseErrorWithMsg(c, CodeInvalidParams, err.Error())
	}

	// 2. 用户登录逻辑
	user, err := logic.Login(&fo)
	if err != nil {
		ResponseError(c, CodeInvalidPassword)
		return
	}

	// 3. 生成token
	accessToken, refreshToken, _ := jwt.GenToken(user.UserID, user.UserName)

	// 4. 返回响应
	ResponseSuccess(c, gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
		"userID":       fmt.Sprintf("%d", user.UserID),
		"username":     user.UserName,
	})
}

// RefreshTokenHandler 刷新accessToken
func RefreshTokenHandler(c *gin.Context) {
	rt := c.Query("refresh_token")
	// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
	// 这里假设Token放在Header的 Authorization 中，并使用 Bearer 开头
	// 这里的具体实现方式要依据你的实际业务情况决定
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		ResponseErrorWithMsg(c, CodeInvalidToken, "请求头缺少Auth Token")
		c.Abort()
		return
	}

	// 按空格分割
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		ResponseErrorWithMsg(c, CodeInvalidToken, "Token格式错误")
		c.Abort()
		return
	}
	aToken, rToken, err := jwt.RefreshToken(parts[1], rt)
	zap.L().Error("jwt.RefreshToken failed", zap.Error(err))
	c.JSON(http.StatusOK, gin.H{
		"access_token":  aToken,
		"refresh_token": rToken,
	})
}
