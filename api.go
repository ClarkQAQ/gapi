package gapi

import (
	"net/http"
	"net/url"
)

// Gapi 接口
type GapiApi struct {
	Method     string
	URL        string
	Headers    http.Header
	Values     url.Values
	Body       []byte
	Error      error
	Hijack     func(p *Gapi, req *http.Request) error
	RespHijack func(resp *GapiResponse, setBody func(b []byte)) error
}

// 新建Gapi接口
func NewAPI(method string, urlString string) *GapiApi {
	a := &GapiApi{
		Method:  method,
		URL:     urlString,
		Headers: http.Header{},
		Values:  url.Values{},
	}

	return a
}

// 设置请求前的Hijack
func (a *GapiApi) SetHijack(hijack func(p *Gapi, req *http.Request) error) *GapiApi {
	a.Hijack = hijack
	return a
}

// 设置请求后的Hijack
func (a *GapiApi) SetRespHijack(hijack func(resp *GapiResponse, setBody func(body []byte)) error) *GapiApi {
	a.RespHijack = hijack
	return a
}

// Headers
// 添加Header
func (a *GapiApi) AddHeader(key string, value string) *GapiApi {
	a.Headers.Add(key, value)
	return a
}

// Headers
// 设置Header
func (a *GapiApi) SetHeader(key string, value string) *GapiApi {
	a.Headers.Set(key, value)
	return a
}

// Headers
// 设置Headers
func (a *GapiApi) SetHeaders(value map[string]string) *GapiApi {
	for k, v := range value {
		a.SetHeader(k, v)
	}
	return a
}

// Headers
// 获取Header
func (a *GapiApi) GetHeader(key string) string {
	return a.Headers.Get(key)
}

// Headers
// 获取Headers
func (a *GapiApi) GetHeaders(key string) []string {
	return a.Headers.Values(key)
}

// Headers
// 删除Header
func (a *GapiApi) DelHeader(key string) {
	a.Headers.Del(key)
}

// Values (URL Query)
// 添加Query
func (a *GapiApi) AddValue(key string, value string) *GapiApi {
	a.Values.Set(key, value)
	return a
}

// Values (URL Query)
// 设置Query
func (a *GapiApi) SetValue(key string, value string) *GapiApi {
	a.Values.Set(key, value)
	return a
}

// Values (URL Query)
// 设置Querys
func (a *GapiApi) SetValues(values map[string]string) *GapiApi {
	for k, v := range values {
		a.Values.Set(k, v)
	}
	return a
}

// Values (URL Query)
// 获取Query
func (a *GapiApi) GetValue(key string) string {
	return a.Values.Get(key)
}

// Values (URL Query)
// 是否存在Query
func (a *GapiApi) HasValue(key string) bool {
	return a.Values.Has(key)
}

// Values (URL Query)
// 删除Query
func (a *GapiApi) DelValue(key string) {
	a.Values.Del(key)
}

// Body
// 设置请求内容
func (a *GapiApi) SetBody(body []byte) *GapiApi {
	a.Body = body
	return a
}

// Body
// 获取请求内容
func (a *GapiApi) GetBody() []byte {
	return a.Body
}

// ERROR
// 抛出错误
func (a *GapiApi) SetError(e error) *GapiApi {
	a.Error = e
	return a
}

// ERROR
// 获取错误
func (a *GapiApi) GetError() error {
	return a.Error
}
