package mysql

import "errors"

// 定义业务状态码

var (
	ErrorUserExist     = errors.New("用户已存在")
	ErrorUserNotExist  = errors.New("用户不已存在")
	ErrorPasswordWrong = errors.New("密码错误")
	ErrorInvalidID     = errors.New("无效的ID")
	ErrorQueryFailed   = errors.New("查询数据失败")
	ErrorInsertFailed  = errors.New("插入数据失败")
)
