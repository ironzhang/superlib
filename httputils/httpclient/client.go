package httpclient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/ironzhang/superlib/codes"
	"github.com/ironzhang/superlib/httputils"
	"github.com/ironzhang/superlib/httputils/httpclient/resolver"
	_ "github.com/ironzhang/superlib/httputils/httpclient/resolver/passthrough"
)

type InvokeInfo struct {
	Method string
	Addr   string
	Path   string
	Query  url.Values
	Header http.Header
	Cookie *http.Cookie
}

type Invoker func(ctx context.Context, info *InvokeInfo, args, reply interface{}) error

type Interceptor func(ctx context.Context, info *InvokeInfo, args, reply interface{}, invoker Invoker) error

type ResultParser func(ctx context.Context, c *Client, resp *http.Response, reply interface{}) error

func makeInterceptorInvoker(i Interceptor, invoker Invoker) Invoker {
	return func(ctx context.Context, info *InvokeInfo, args, reply interface{}) error {
		return i(ctx, info, args, reply, invoker)
	}
}

type Client struct {
	Addr         string
	Codec        Codec
	Client       http.Client
	ParseResult  ResultParser
	Interceptors []Interceptor

	once    sync.Once
	target  resolver.Target
	invoker Invoker
}

func makeTarget(addr string) resolver.Target {
	results := strings.SplitN(addr, "/", 2)
	if len(results) == 2 {
		return resolver.Target{
			Scheme:   results[0],
			Endpoint: results[1],
		}
	}
	return resolver.Target{
		Endpoint: addr,
	}
}

func (p *Client) init() {
	// 设置默认编解码器
	if p.Codec == nil {
		p.Codec = JSONCodec{}
	}

	// 设置默认结果解析函数
	if p.ParseResult == nil {
		p.ParseResult = defaultResultParser
	}

	// 构造目标地址
	p.target = makeTarget(p.Addr)

	// 构造拦截器调用链
	invoker := p.call
	for i := len(p.Interceptors) - 1; i >= 0; i-- {
		invoker = makeInterceptorInvoker(p.Interceptors[i], invoker)
	}
	p.invoker = invoker
}

func (p *Client) NewAddrClient(addr string) *Client {
	c := p.clone()
	c.Addr = addr
	return c
}

func (p *Client) NewCodecClient(codec Codec) *Client {
	c := p.clone()
	c.Codec = codec
	return c
}

func (p *Client) clone() *Client {
	return &Client{
		Addr:         p.Addr,
		Codec:        p.Codec,
		Client:       p.Client,
		ParseResult:  p.ParseResult,
		Interceptors: p.Interceptors,
	}
}

func (p *Client) Get(ctx context.Context, path string, query url.Values, reply interface{}, opts ...CallOption) error {
	return p.Invoke(ctx, "GET", path, query, nil, reply, opts...)
}

func (p *Client) Post(ctx context.Context, path string, query url.Values, args, reply interface{}, opts ...CallOption) error {
	return p.Invoke(ctx, "POST", path, query, args, reply, opts...)
}

func (p *Client) Invoke(ctx context.Context, method, path string, query url.Values, args, reply interface{}, opts ...CallOption) error {
	// 初始化
	p.once.Do(p.init)

	// 解析地址
	addr, err := resolver.Resolve(p.target)
	if err != nil {
		return fmt.Errorf("resolve: %w", err)
	}

	// 构造调用信息
	if query == nil {
		query = make(url.Values)
	}
	info := InvokeInfo{
		Method: method,
		Addr:   addr,
		Path:   path,
		Query:  query,
		Header: make(http.Header),
	}
	for _, o := range opts {
		o(&info)
	}

	// 执行调用链
	return p.invoker(ctx, &info, args, reply)
}

func normalizePath(path string) string {
	if path == "" {
		path = "/"
	} else if path[0] != '/' {
		path = "/" + path
	}
	return path
}

func (p *Client) call(ctx context.Context, info *InvokeInfo, args, reply interface{}) error {
	var err error
	var body bytes.Buffer

	// 序列化消息体
	if args != nil {
		if err = p.Codec.Encode(&body, args); err != nil {
			return fmt.Errorf("codec encode: %w", err)
		}
	}

	// 构造 url
	url := "http://" + info.Addr + normalizePath(info.Path)
	if len(info.Query) > 0 {
		url = url + "?" + info.Query.Encode()
	}

	// 构造请求
	req, err := http.NewRequestWithContext(ctx, info.Method, url, &body)
	if err != nil {
		return fmt.Errorf("new http request with context: %w", err)
	}
	req.Header = info.Header.Clone()
	if args != nil {
		req.Header.Set(httputils.HeaderContentType, p.Codec.ContentType())
	}
	if info.Cookie != nil {
		req.AddCookie(info.Cookie)
	}

	// 发送请求
	resp, err := p.Client.Do(req)
	if err != nil {
		return fmt.Errorf("http client do: %w", err)
	}
	defer resp.Body.Close()

	// 解析响应
	return p.ParseResult(ctx, p, resp, reply)
}

func defaultResultParser(ctx context.Context, c *Client, resp *http.Response, reply interface{}) (err error) {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		data, err := read(resp.Body, resp.ContentLength)
		if err != nil {
			return fmt.Errorf("codec read body: %w [%s]", err, resp.Status)
		}
		buf := bytes.NewBuffer(data)

		var m codes.Message
		if err = c.Codec.Decode(buf, &m); err != nil {
			return fmt.Errorf("codec decode: %w [%s] %s", err, resp.Status, data)
		}
		return codes.MessageError(m)
	}
	if reply != nil {
		if err = c.Codec.Decode(resp.Body, reply); err != nil {
			return fmt.Errorf("codec decode: %w [%s]", err, resp.Status)
		}
	}
	return nil
}

func read(r io.Reader, length int64) ([]byte, error) {
	if length > 0 {
		data := make([]byte, length)
		_, err := io.ReadFull(r, data)
		return data, err
	}
	return io.ReadAll(r)
}
