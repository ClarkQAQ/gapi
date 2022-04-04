package api

import (
	"net/http"
	"net/url"
)

// Pixiv 接口
type PixivApi struct {
	Method  string
	Path    string
	Headers http.Header
	Values  url.Values
	Body    []byte
	Error   error
}

// 新建Pixiv接口
func NewApi(method string, path string, headers http.Header, values url.Values, body []byte, e error) *PixivApi {
	return &PixivApi{
		Method:  method,
		Path:    path,
		Headers: headers,
		Values:  values,
		Body:    body,
		Error:   e,
	}
}
