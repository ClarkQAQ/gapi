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
func New(method string, urlString string, headers http.Header, values url.Values, body []byte, e error) *PixivApi {
	a := &PixivApi{
		Method:  method,
		URL:     urlString,
		Headers: headers,
		Values:  values,
		Body:    body,
		Error:   e,
	}

	return a
}

func (a *PixivApi) SetHijack(hijack func(c *http.Client, req *http.Request) error) {
	a.Hijack = hijack
}

func (a *PixivApi) SetRespHijack(hijack func(resp *http.Response, respBody func(b []byte) []byte) error) {
	a.RespHijack = hijack
}
