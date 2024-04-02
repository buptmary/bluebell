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

// GetPostList2Handler 升级版帖子列表接口
func GetPostList2Handler(c *gin.Context) {
	// 1. 获取参数
	// GET请求参数(query string)：/api/v1/posts2?page=1&size=10&order=time
	// 初始化结构体时指定初始参数
	p := &models.PostListForm{
		Page:  1,
		Size:  10,
		Order: models.OrderTime,
	}
	//c.ShouldBind() 根据请求的数据类型选择相应的方法去获取数据
	//c.ShouldBindJSON() 如果请求中携带的是json格式的数据，才能用这个方法获取到数据
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostList2Handler with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 2. 去Redis查询id列表
	// 3. 根据id去数据库查询帖子详细信息
	data, err := logic.GetPostList2(p)
	if err != nil {
		zap.L().Error("logic.GetPostList2(p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 4. 成功返回贴子数据
	ResponseSuccess(c, data)

}
