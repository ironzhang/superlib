package httputils

import (
	"github.com/labstack/echo"
)

// 公共头部名称定义
const (
	// 通用头部名称
	HeaderContentType = "Content-Type"

	// 自定义头部名称
	HeaderXRequestID         = echo.HeaderXRequestID
	HeaderXTraceID           = "X-Trace-Id"
	HeaderXSpanID            = "X-Span-Id"
	HeaderXTaskID            = "X-Task-Id"
	HeaderXCaller            = "X-Caller"
	HeaderAccessLogPrintMode = "X-Access-Log-Print-Mode"
)

// 公共头部值定义
const (
	ApplicationJSON = "application/json"
	ApplicationForm = "application/x-www-form-urlencoded"
)

// AccessLogPrintMode 访问日志输出模式
type AccessLogPrintMode string

const (
	DisableAccessBody AccessLogPrintMode = "DisableAccessBody"
)
