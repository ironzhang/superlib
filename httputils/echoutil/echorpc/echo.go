package echorpc

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/labstack/echo"

	"github.com/ironzhang/superlib/codes"
)

func HandlerFunc(i interface{}) echo.HandlerFunc {
	f, e := parseFunction(i)
	return func(c echo.Context) error {
		if e != nil {
			return fmt.Errorf("parse function: %w", e)
		}

		var err error
		args := f.NewArgs()
		reply := f.NewReply()

		// decode
		if !isNilInterface(f.args) {
			if err = c.Bind(args.Interface()); err != nil {
				return codes.Errorf(codes.InvalidParams, "bind: %w", err)
			}
		}
		if f.args.Kind() != reflect.Ptr {
			args = args.Elem()
		}

		// validate
		if v := c.Echo().Validator; v != nil {
			if err = v.Validate(args.Interface()); err != nil {
				return codes.Errorf(codes.InvalidParams, "validate: %w", err)
			}
		}

		// call
		if err = f.Call(c.Request().Context(), args, reply); err != nil {
			return err
		}

		// reply
		if !isNilInterface(f.reply) {
			return c.JSON(http.StatusOK, reply.Interface())
		}
		msg := codes.Message{
			Code: int(codes.OK),
			Desc: codes.OK.String(),
		}
		return c.JSON(http.StatusOK, msg)
	}
}
