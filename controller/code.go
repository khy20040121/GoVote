package controller

type ResCode int64

const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParam
	CodeUserExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeVoteRepeated
	CodeServerBusy
	CodePostNotExist

	CodeNeedLogin
	CodeInvalidToken
)

var codeMsg = map[ResCode]string{
	CodeSuccess:         "success",
	CodeInvalidParam:    "请求参数有误",
	CodeUserExist:       "用户已存在",
	CodeUserNotExist:    "用户不存在",
	CodeInvalidPassword: "用户名或密码错误",
	CodeServerBusy:      "服务器繁忙",
	CodeVoteRepeated:    "用户重复投票",
	CodeNeedLogin:       "用户未登录",
	CodeInvalidToken:    "Token已失效",
	CodePostNotExist:    "查询不到帖子",
}

func (c ResCode) Msg() string {
	return codeMsg[c]
}
