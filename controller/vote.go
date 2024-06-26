package controller

import (
	"bluebell/dao/redis"
	"bluebell/logic"
	"bluebell/models"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// PostVoteHandler 投票函数
func PostVoteHandler(c *gin.Context) {
	// 1. 获取参数及参数校验
	vote := new(models.VoteDataForm)
	if err := c.ShouldBindJSON(vote); err != nil {
		var errs validator.ValidationErrors
		ok := errors.As(err, &errs) // 类型断言
		if !ok {
			ResponseError(c, CodeInvalidParams)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParams, "投票参数错误")
		return
	}
	// 获取当前用户的id
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNotLogin)
		return
	}

	// 具体投票的业务逻辑
	if err := logic.VoteForPost(userID, vote); err != nil {
		zap.L().Error("logic.VoteForPost() failed", zap.Error(err))
		switch {
		case errors.Is(err, redis.ErrorVoteRepeated):
			ResponseError(c, CodeVoteRepeated)
		case errors.Is(err, redis.ErrorVoteTimeExpire):
			ResponseError(c, CodeVoteTimeExpire)
		default:
			ResponseError(c, CodeServerBusy)
		}
		return
	}

	// 返回成功响应
	ResponseSuccess(c, nil)
}
