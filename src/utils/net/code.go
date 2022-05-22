package net

type ResCode int64

// CODE
const (
	// 成功（默认返回状态码）
	CodeSuccess ResCode = 0
	// 全局未知异常
	CodeSeverError ResCode = 500
	// 请求失败（一般前端处理，不常用）
	CodeBadRequest ResCode = 400
	// 请求资源不存在（静态资源不存在，不常用）
	CodeDataNotFount ResCode = 404
	// 登录、权限认证异常
	CodeLoginExpire ResCode = 401
	// 权限不足
	CodeIdentityNotRow ResCode = 403
)

// CODE_MEAN
var codeMsgMap = map[ResCode]string{
	CodeSuccess:        "success",
	CodeSeverError:     "服务器繁忙请重试",
	CodeBadRequest:     "请求失败",
	CodeDataNotFount:   "未找到资源",
	CodeLoginExpire:    "请登录后重试",
	CodeIdentityNotRow: "权限不足",
}

// 获取CODE意义
func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeSeverError]
	}
	return msg
}
