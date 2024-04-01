package controller

// ResCode 定义返回的状态码类型
type ResCode int64

const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParams
	CodeUserExist
	CodeUserNotExist
	CodeInvalidPassword

	CodeServerBusy
	CodeInvalidToken
	CodeInvalidAuthFormat
	CodeNotLogin
	CodeVoteRepeated
	CodeVoteTimeExpire
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:         "success",
	CodeInvalidParams:   "请求参数错误",
	CodeUserExist:       "用户已存在",
	CodeUserNotExist:    "用户不存在",
	CodeInvalidPassword: "用户名或密码错误",

	CodeServerBusy:        "服务繁忙",
	CodeInvalidToken:      "无效的Token",
	CodeInvalidAuthFormat: "认证格式有误",
	CodeNotLogin:          "未登录",
	CodeVoteRepeated:      "请勿重复投票",
	CodeVoteTimeExpire:    "投票时间已过",
}

func (code ResCode) Msg() string {
	msg, ok := codeMsgMap[code] // 从字典查询当前状态码信息
	if ok {
		return msg
	}
	return codeMsgMap[CodeServerBusy]
}
