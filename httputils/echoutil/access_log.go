package echoutil

import (
	"bytes"
	"net/http"
	"time"

	"github.com/labstack/echo"

	"github.com/ironzhang/tlog"

	"github.com/ironzhang/superlib/codes"
	"github.com/ironzhang/superlib/httputils"
)

type DisableBodyFunc func(c echo.Context) bool

// DisableBodyByPath 指定的path不进行body的打印
func DisableBodyByPath(paths []string) DisableBodyFunc {
	disablePathMap := make(map[string]struct{})
	for _, path := range paths {
		disablePathMap[path] = struct{}{}
	}
	return func(c echo.Context) bool {
		path := c.Request().URL.Path
		if _, ok := disablePathMap[path]; ok {
			return true
		}
		return false
	}
}

// AccessLogConfig 访问日志中间件配置
type AccessLogConfig struct {
	DisableBody DisableBodyFunc
}

func (p *AccessLogConfig) MiddlewareFunc(next echo.HandlerFunc) echo.HandlerFunc {
	log := tlog.Named("access.in")
	return func(c echo.Context) error {
		// 判断日志是否需要打印 body
		disableBody := p.isDisableAccessLogBody(c)

		// 输出请求日志
		p.printRequest(c, log, disableBody)

		// 构建响应记录器
		resp := c.Response()
		recorder := httputils.NewResponseRecorder(resp.Writer)
		resp.Writer = recorder

		// 调用下一个中间件
		start := time.Now()
		err := next(c)

		// 输出响应日志
		p.printResponse(c, log, time.Since(start), err, recorder, disableBody)

		// 返回错误
		return err
	}
}

func (p *AccessLogConfig) printRequest(c echo.Context, log tlog.Logger, disableBody bool) {
	req := c.Request()
	ctx := req.Context()

	// 输出日志
	log = log.WithContext(ctx).WithArgs("method", req.Method, "path", req.URL.Path,
		"query", req.URL.Query(), "header", req.Header, "remote_addr", req.RemoteAddr)
	if disableBody {
		log.Infof("http server request")
	} else {
		var data []byte
		var err error
		req.Body, data, err = httputils.CopyBody(req.Body)
		if err != nil {
			log.Infow("copy body", "error", err)
			return
		}
		// 输出 body 日志
		log.Infof("http server request: %s", trimNewline(data))
	}
}

func (p *AccessLogConfig) printResponse(c echo.Context, log tlog.Logger, latency time.Duration, err error, r *httputils.ResponseRecorder, disableBody bool) {
	req := c.Request()
	ctx := req.Context()

	// 输出日志
	log = log.WithContext(ctx).WithArgs("method", req.Method, "path", req.URL.Path,
		"latency", latency, "error", err, "code", codes.GetErrorCode(err), "status", http.StatusText(r.Status()), "header", r.Header(), "remote_addr", req.RemoteAddr)
	if disableBody {
		log.Infof("http server response")
	} else {
		log.Infof("http server response: %s", trimNewline(r.Body()))
	}
}

func trimNewline(data []byte) []byte {
	return bytes.TrimSuffix(data, []byte{'\n'})
}

// AccessLogMiddleware 访问日志中间件
func AccessLogMiddleware() echo.MiddlewareFunc {
	cfg := AccessLogConfig{}
	return cfg.MiddlewareFunc
}

func (p *AccessLogConfig) isDisableAccessLogBody(c echo.Context) bool {
	disableBody := false
	if p.DisableBody != nil {
		disableBody = p.DisableBody(c)
	}
	if disableBody {
		return true
	}
	return c.Request().Header.Get(httputils.HeaderAccessLogPrintMode) == string(httputils.DisableAccessBody)
}
