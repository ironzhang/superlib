package codes

import (
	"fmt"
	"net/http"
)

type Code int

var codes = map[Code]string{}

func Register(code Code, desc string) {
	registered, ok := codes[code]
	if ok {
		panic(fmt.Sprintf("code=%d(%s,%s) is registered", code, desc, registered))
	}
	codes[code] = desc
}

func (c Code) String() string {
	if desc, ok := codes[c]; ok {
		return desc
	}
	return fmt.Sprintf("code(%d)", c)
}

const (
	// 0~999 公共错误码
	OK               Code = 0 // 成功
	Unknown          Code = 1 // 未知错误
	Internal         Code = 2 // 内部错误
	NotFound         Code = 3 // 未找到对象
	InvalidParams    Code = 4 // 无效参数
	DeadlineExceeded Code = 5 // 调用超时
	Unavailable      Code = 6 // 服务不可用
	Canceled         Code = 7 // 取消
	Unauthorized     Code = 8 // 未授权
	Forbidden        Code = 9 // 禁止访问
)

func init() {
	Register(OK, "ok")
	Register(Unknown, "unknown")
	Register(Internal, "internal")
	Register(NotFound, "not found")
	Register(InvalidParams, "invalid parameters")
	Register(DeadlineExceeded, "deadline exceeded")
	Register(Unavailable, "unavailable")
	Register(Canceled, "context canceled")
	Register(Unauthorized, "unauthorized")
	Register(Forbidden, "forbidden")
}

func HTTPStatus(code Code) int {
	switch code {
	case OK:
		return http.StatusOK
	case InvalidParams:
		return http.StatusBadRequest
	case Unauthorized:
		return http.StatusUnauthorized
	case Forbidden:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
