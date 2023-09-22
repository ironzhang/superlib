package httputils

import (
	"bytes"
	"net/http"
)

type ResponseRecorder struct {
	http.ResponseWriter
	status int
	buffer bytes.Buffer
}

func NewResponseRecorder(w http.ResponseWriter) *ResponseRecorder {
	return &ResponseRecorder{ResponseWriter: w}
}

func (p *ResponseRecorder) WriteHeader(status int) {
	p.status = status
	p.ResponseWriter.WriteHeader(status)
}

func (p *ResponseRecorder) Write(b []byte) (int, error) {
	p.buffer.Write(b)
	return p.ResponseWriter.Write(b)
}

func (p *ResponseRecorder) Status() int {
	return p.status
}

func (p *ResponseRecorder) Body() []byte {
	return p.buffer.Bytes()
}
