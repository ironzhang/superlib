package interceptors

import "github.com/ironzhang/superlib/codes"

func getErrorCode(err error) codes.Code {
	return codes.GetErrorCode(err)
}
