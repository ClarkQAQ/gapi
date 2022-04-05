package api

import (
	"net/http"
	"net/url"
)

// Pixiv 接口
type PixivApi struct {
	Method     string
	URL        string
	Headers    http.Header
	Values     url.Values
	Body       []byte
	Error      error
	Hijack     func(c *http.Client, req *http.Request) error
	RespHijack func(resp *http.Response, respBody func(b []byte) []byte) error
}

// 新建Pixiv接口
func New(method string, urlString string) *PixivApi {
	a := &PixivApi{
		Method:  method,
		URL:     urlString,
		Headers: http.Header{},
		Values:  url.Values{},
	}

	return a
}

func (a *PixivApi) SetHijack(hijack func(c *http.Client, req *http.Request) error) *PixivApi {
	a.Hijack = hijack
	return a
}

func (a *PixivApi) SetRespHijack(hijack func(resp *http.Response, respBody func(b []byte) []byte) error) *PixivApi {
	a.RespHijack = hijack
	return a
}

func (a *PixivApi) AddHeader(key string, value string) *PixivApi {
	a.Headers.Add(key, value)
	return a
}

func (a *PixivApi) SetHeader(key string, value string) *PixivApi {
	a.Headers.Set(key, value)
	return a
}

func (a *PixivApi) GetHeader(key string) string {
	return a.Headers.Get(key)
}

func (a *PixivApi) GetHeaders(key string) []string {
	return a.Headers.Values(key)
}

func (a *PixivApi) DelHeader(key string) {
	a.Headers.Del(key)
}

func (a *PixivApi) AddValue(key string, value string) *PixivApi {
	a.Values.Set(key, value)
	return a
}

func (a *PixivApi) SetValue(key string, value string) *PixivApi {
	a.Values.Set(key, value)
	return a
}

func (a *PixivApi) GetValue(key string) string {
	return a.Values.Get(key)
}

func (a *PixivApi) HasValue(key string) bool {
	return a.Values.Has(key)
}

func (a *PixivApi) DelValue(key string) {
	a.Values.Del(key)
}

func (a *PixivApi) SetBody(body []byte) *PixivApi {
	a.Body = body
	return a
}

func (a *PixivApi) GetBody() []byte {
	return a.Body
}

func (a *PixivApi) SetError(e error) *PixivApi {
	a.Error = e
	return a
}

func (a *PixivApi) GetError() error {
	return a.Error
}
