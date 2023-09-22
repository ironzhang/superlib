package echoutil

import (
	"net/http"

	"github.com/labstack/echo"
)

// ServeMuxConfig 多路复用中间件配置
type ServeMuxConfig struct {
	ServeMux *http.ServeMux
}

func (p *ServeMuxConfig) MiddlewareFunc(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		h, pattern := p.ServeMux.Handler(req)
		if pattern != "" {
			h.ServeHTTP(c.Response(), req)
			return nil
		}
		return next(c)
	}
}

// ServeMuxMiddleware 多路复用中间件
func ServeMuxMiddleware(m *http.ServeMux) echo.MiddlewareFunc {
	if m == nil {
		m = http.DefaultServeMux
	}
	cfg := ServeMuxConfig{
		ServeMux: m,
	}
	return cfg.MiddlewareFunc
}
