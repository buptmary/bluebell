package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
)

var (
	ErrorUserNotLogin = errors.New("当前用户未登录")
)

func getCurrentUserID(c *gin.Context) (userID int64, err error) {
	_userID, ok := c.Get(ContextUserIDKey)
	if !ok {
		err = ErrorUserNotLogin
		return 0, err
	}

	userID, ok = _userID.(int64) // 类型断言
	if !ok {
		err = ErrorUserNotLogin
		return 0, err
	}
	return userID, nil
}

// getPageInfo 分页参数
func getPageInfo(c *gin.Context) (int64, int64) {
	pageStr := c.Query("page")
	sizeStr := c.Query("size")

	var (
		page int64 // 第几页 页数
		size int64 // 每页size条数据
		err  error
	)
	page, err = strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		page = 1
	}
	size, err = strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		size = 1
	}

	return page, size
}
