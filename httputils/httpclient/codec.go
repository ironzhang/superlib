package httpclient

import (
	"encoding/json"
	"io"
	"net/url"

	"github.com/ironzhang/superlib/httputils"
	"github.com/ironzhang/superlib/httputils/bind"
)

type Codec interface {
	ContentType() string
	Encode(w io.Writer, v interface{}) error
	Decode(r io.Reader, v interface{}) error
}

type JSONCodec struct {
}

func (p JSONCodec) ContentType() string {
	return httputils.ApplicationJSON
}

func (p JSONCodec) Encode(w io.Writer, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

func (p JSONCodec) Decode(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

// FormJSONCodec 请求按照form进行编码，响应按照json进行解码
type FormJSONCodec struct{}

func (FormJSONCodec) ContentType() string {
	return httputils.ApplicationForm
}

func (FormJSONCodec) Encode(w io.Writer, v interface{}) error {
	if values, ok := v.(url.Values); ok {
		io.WriteString(w, values.Encode())
		return nil
	}

	if data, ok := v.(map[string]string); ok {
		values := url.Values{}
		for k, v := range data {
			values.Add(k, v)
		}
		io.WriteString(w, values.Encode())
		return nil
	}

	values, err := bind.BindQuery(v, "json")
	if err != nil {
		return err
	}
	io.WriteString(w, values.Encode())
	return nil
}

func (FormJSONCodec) Decode(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}
