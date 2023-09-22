package echoutil

import (
	"github.com/labstack/echo"

	"github.com/ironzhang/superlib/codes"
)

func HTTPErrorHandler(err error, c echo.Context) {
	if err == nil {
		return
	}

	msg := codes.ErrorMessage(err)
	status := codes.HTTPStatus(codes.Code(msg.Code))
	if he, ok := err.(*echo.HTTPError); ok {
		status = he.Code
	}
	c.JSON(status, msg)
}
