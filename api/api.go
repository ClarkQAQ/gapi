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

// 设置请求前的Hijack
func (a *PixivApi) SetHijack(hijack func(c *http.Client, req *http.Request) error) *PixivApi {
	a.Hijack = hijack
	return a
}

// 设置请求后的Hijack
func (a *PixivApi) SetRespHijack(hijack func(resp *http.Response, respBody func(b []byte) []byte) error) *PixivApi {
	a.RespHijack = hijack
	return a
}

// Headers
// 添加Header
func (a *PixivApi) AddHeader(key string, value string) *PixivApi {
	a.Headers.Add(key, value)
	return a
}

// Headers
// 设置Header
func (a *PixivApi) SetHeader(key string, value string) *PixivApi {
	a.Headers.Set(key, value)
	return a
}

// Headers
// 获取Header
func (a *PixivApi) GetHeader(key string) string {
	return a.Headers.Get(key)
}

// Headers
// 获取Headers
func (a *PixivApi) GetHeaders(key string) []string {
	return a.Headers.Values(key)
}

// Headers
// 删除Header
func (a *PixivApi) DelHeader(key string) {
	a.Headers.Del(key)
}

// Values (URL Query)
// 添加Query
func (a *PixivApi) AddValue(key string, value string) *PixivApi {
	a.Values.Set(key, value)
	return a
}

// Values (URL Query)
// 设置Query
func (a *PixivApi) SetValue(key string, value string) *PixivApi {
	a.Values.Set(key, value)
	return a
}

// Values (URL Query)
// 获取Query
func (a *PixivApi) GetValue(key string) string {
	return a.Values.Get(key)
}

// Values (URL Query)
// 是否存在Query
func (a *PixivApi) HasValue(key string) bool {
	return a.Values.Has(key)
}

// Values (URL Query)
// 删除Query
func (a *PixivApi) DelValue(key string) {
	a.Values.Del(key)
}

// Body
// 设置请求内容
func (a *PixivApi) SetBody(body []byte) *PixivApi {
	a.Body = body
	return a
}

// Body
// 获取请求内容
func (a *PixivApi) GetBody() []byte {
	return a.Body
}

// ERROR
// 抛出错误
func (a *PixivApi) SetError(e error) *PixivApi {
	a.Error = e
	return a
}

// ERROR
// 获取错误
func (a *PixivApi) GetError() error {
	return a.Error
}
