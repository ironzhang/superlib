package httputils

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

// CopyBody 从 body 中复制数据
func CopyBody(b io.ReadCloser) (r io.ReadCloser, data []byte, err error) {
	if b == nil || b == http.NoBody {
		return http.NoBody, nil, nil
	}
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return b, nil, err
	}
	if err = b.Close(); err != nil {
		return b, nil, err
	}
	return ioutil.NopCloser(&buf), buf.Bytes(), nil
}
