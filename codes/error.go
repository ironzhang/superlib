package codes

import (
	"context"
	"errors"
	"fmt"
)

type Message struct {
	Code  int    ``                  // 错误码
	Desc  string ``                  // 错误码描述
	Error string `json:",omitempty"` // 错误详情
}

type Error struct {
	code Code
	wrap error
}

func NewError(code Code, err error) Error {
	return Error{
		code: code,
		wrap: err,
	}
}

func Errorf(code Code, format string, a ...interface{}) Error {
	return Error{
		code: code,
		wrap: fmt.Errorf(format, a...),
	}
}

func (e Error) Code() Code {
	return e.code
}

func (e Error) Unwrap() error {
	return e.wrap
}

func (e Error) Error() string {
	return fmt.Sprintf("{code: %d, desc: %s, wrap: %v}", e.code, e.code.String(), e.wrap)
}

func GetErrorCode(err error) Code {
	if err == nil {
		return OK
	}
	c, ok := err.(interface {
		Code() Code
	})
	if ok {
		return c.Code()
	}
	w := errors.Unwrap(err)
	if w == nil {
		return parseErrorCode(err)
	}
	return GetErrorCode(w)
}

func parseErrorCode(err error) Code {
	switch err {
	case context.Canceled:
		return Canceled
	case context.DeadlineExceeded:
		return DeadlineExceeded
	default:
		if isCanceledErr(err) {
			return Canceled
		}
		if isTimeoutErr(err) {
			return DeadlineExceeded
		}
		return Unknown
	}
}

func isCanceledErr(err error) bool {
	// 兼容处理 net 包的中的 errCanceled
	if err.Error() == "operation was canceled" {
		return true
	}
	return false
}

func isTimeoutErr(err error) bool {
	// 兼容处理 Timeout 接口，如 net 包中的 errTimeout
	t, ok := err.(interface {
		Timeout() bool
	})
	if ok {
		return t.Timeout()
	}
	return false
}

func ErrorMessage(err error) Message {
	code := GetErrorCode(err)
	return Message{
		Code:  int(code),
		Desc:  code.String(),
		Error: err.Error(),
	}
}

func MessageError(m Message) Error {
	return Errorf(Code(m.Code), m.Error)
}
