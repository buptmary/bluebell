package controller

import (
	"bluebell/logic"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

// CommunityDetailHandler 根据ID查找到社区分类的详情
func CommunityDetailHandler(c *gin.Context) {
	// 1. 获取社区id
	communityIdStr := c.Param("id")
	communityId, err := strconv.ParseInt(communityIdStr, 10, 64)
	if err != nil {
		zap.L().Error("community id parse failed", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 2. 根据id获取社区详情
	communityDetail, err := logic.GetCommunityByID(communityId)
	if err != nil {
		zap.L().Error("get community by id failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 成功，返回communityDetail
	ResponseSuccess(c, communityDetail)
}

// CommunityHandler 查找社区列表
func CommunityHandler(c *gin.Context) {
	// 查询到所有的社区 以列表形式返回
	communityList, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList()", zap.Error(err))
		ResponseError(c, CodeServerBusy) // 不轻易将服务端报错返回给客户端
		return
	}

	// 成功返回社区列表
	ResponseSuccess(c, communityList)
}
