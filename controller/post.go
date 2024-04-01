package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

func CreatePostHandler(c *gin.Context) {
	// 1. 获取参数及参数校验
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil { // validator --> binding tag
		zap.L().Error("c.ShouldBindJSON(&post)", zap.Error(err))
		ResponseErrorWithMsg(c, CodeInvalidParams, err.Error())
		return
	}
	// 2. 获取用户ID (当前请求的UserID)
	userID, err := getCurrentUserID(c)
	if err != nil {
		zap.L().Error("getCurrentUserID(c) failed", zap.Error(err))
		ResponseError(c, CodeNotLogin)
		return
	}
	post.AuthorID = userID

	// 3. 创建帖子

	if err := logic.CreatePost(&post); err != nil {
		zap.L().Error("logic.CreatePost(&post) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 4. 返回响应
	ResponseSuccess(c, nil)
}

// PostDetailHandler 根据Id查询帖子详情
func PostDetailHandler(c *gin.Context) {
	// 1. 获取请求参数及参数校验
	postIdStr := c.Param("id")
	postId, err := strconv.ParseInt(postIdStr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 2. 根据id取出对应帖子数据
	post, err := logic.GetPostByID(postId)
	if err != nil {
		zap.L().Error("logic.GetPostByID(postId) failed", zap.Error(err))
		ResponseErrorWithMsg(c, CodeInvalidParams, err.Error())
		return
	}

	// 3. 成功，返回响应
	ResponseSuccess(c, post)
}

// PostListHandler 分页获取帖子列表
func PostListHandler(c *gin.Context) {
	// 获取分页参数
	page, size := getPageInfo(c)
	// 获取分页数据
	data, err := logic.GetPostList(page, size)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}
