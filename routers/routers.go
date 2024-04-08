package routers

import (
	"bluebell/controller"
	"bluebell/middlewares"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middlewares.RateLimitMiddleware(2*time.Second, 10)) // 每2秒添加1个令牌 容量为40
	v1 := r.Group("/api/v1")
	// 用户注册
	v1.POST("/login", controller.LoginHandler)
	// 用户登录
	v1.POST("/signup", controller.SignUpHandler)
	// 刷新Token
	v1.GET("/refresh_token", controller.RefreshTokenHandler)

	v1.Use(controller.JWTAuth())
	{
		v1.GET("/community/:id", controller.CommunityDetailHandler)
		v1.GET("/community", controller.CommunityHandler)

		v1.POST("/post", controller.CreatePostHandler)
		v1.GET("/post/:id", controller.PostDetailHandler)
		v1.GET("/posts", controller.PostListHandler)
		// 根据时间或分数获取帖子列表
		v1.GET("/posts2", controller.GetPostList2Handler)

		// 投票
		v1.POST("/vote", controller.PostVoteHandler)

		v1.GET("/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "pong")
		})
	}

	// 处理未注册路由
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404 page not found!",
		})
	})

	return r
}
