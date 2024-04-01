package controller

import (
	"bluebell/pkg/jwt"
	"github.com/gin-gonic/gin"
	"strings"
)

const (
	ContextUserIDKey = "userID"
)

func JWTAuth() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
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
		// parts[1]是获取到的tokenString, 使用JWT解析函数进行解析
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			ResponseError(c, CodeInvalidToken)
			c.Abort()
			return
		}

		// 将请求的userid信息保存到上下文c中
		c.Set(ContextUserIDKey, mc.UserID)
		c.Next() // 后续的处理函数可以用c.Get("userID")来获取当前请求的用户信息
	}
}
