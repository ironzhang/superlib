package echoutil

import (
	"fmt"
	"runtime"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/ironzhang/tlog"
)

// Recover recover 中间件
func Recover() echo.MiddlewareFunc {
	return recoverWithConfigAndMetrics(middleware.DefaultRecoverConfig)
}

// recoverWithConfigAndMetrics returns a Recover middleware with config.
func recoverWithConfigAndMetrics(config middleware.RecoverConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = middleware.DefaultRecoverConfig.Skipper
	}
	if config.StackSize == 0 {
		config.StackSize = middleware.DefaultRecoverConfig.StackSize
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					path := c.Request().URL.Path
					stack := make([]byte, config.StackSize)
					length := runtime.Stack(stack, !config.DisableStackAll)
					if !config.DisablePrintStack {
						tlog.Named("panic").Errorf("echo handler panic, path:%s, err:%v, stack:\n%s\n", path, err, stack[:length])
					}
					c.Error(err)

					// metrics上报
					// metrics.Counter("public.panic", nil)
				}
			}()
			return next(c)
		}
	}
}
