package gapi

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"utilware/gjson"
)

// Gapi接口返回
type GapiResponse struct {
	p       *Gapi
	api     *GapiApi
	raw     *bytes.Buffer
	content []byte
	result  *gjson.Result

	*http.Response
}

// 请求Gapi接口 (接口需通过api包里面的方法生成)
func (p *Gapi) Do(api *GapiApi) (presp *GapiResponse, e error) {
	// 抛出上层API的错误
	if api.Error != nil {
		return nil, api.Error
	}

	// 判断是否是URL, 如果是URL就执行URL的方式
	u := p.EndpointPATH(api.URL, api.Values)
	if strings.Contains(api.URL, "://") {
		u, e = p.EndpointURL(api.URL, api.Values)
		if e != nil {
			return nil, e
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	// Body
	var body io.Reader = nil
	if api.Body != nil {
		body = bytes.NewReader(api.Body)
	}

	// 设置请求体
	resp, e := p.Request(ctx, api.Method, u.String(), body, func(c *http.Client, req *http.Request) error {
		// 设置请求头
		if api.Headers != nil && len(api.Headers) > 0 {
			for k, v := range api.Headers {
				for i := 0; i < len(v); i++ {
					req.Header.Set(k, v[i])
				}
			}
		}

		// 执行HiJack
		if api.Hijack != nil {
			if e := api.Hijack(p, req); e != nil {
				return e
			}
		}

		return nil
	})

	if e != nil {
		return nil, e
	}

	presp = &GapiResponse{
		p:        p,
		api:      api,
		Response: resp,
	}

	b, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return nil, e
	}

	presp.raw = bytes.NewBuffer(b)
	defer resp.Body.Close()

	if api.RespHijack != nil {
		if e := api.RespHijack(presp, func(body []byte) {
			presp.raw = bytes.NewBuffer(body)
			presp.content = nil
			presp.result = nil
		}); e != nil {
			return nil, e
		}
	}

	return presp, nil
}

// 返回原始的响应内容
// 可以多次调用
func (r *GapiResponse) Raw() []byte {
	return r.raw.Bytes()
}

// 获取响应内容
// 会自动解压缩
func (r *GapiResponse) Content() ([]byte, error) {
	if r.content != nil {
		return r.content, nil
	}

	rawBytes := r.Raw()
	var reader io.ReadCloser
	var e error

	switch r.Header.Get(HeaderContentEncoding) {
	case "gzip":
		if reader, e = gzip.NewReader(bytes.NewBuffer(r.raw.Bytes())); e != nil {
			return nil, e
		}
	case "deflate":
		// deflate should be zlib
		// http://www.gzip.org/zlib/zlib_faq.html#faq38
		if reader, e = zlib.NewReader(bytes.NewBuffer(r.raw.Bytes())); e != nil {
			// try RFC 1951 deflate
			// http: //www.open-open.com/lib/view/open1460866410410.html
			reader = flate.NewReader(bytes.NewBuffer(r.raw.Bytes()))
		}
	}

	if reader == nil {
		r.content = rawBytes
		return rawBytes, nil
	}

	defer reader.Close()
	b, e := ioutil.ReadAll(reader)

	if e != nil {
		return nil, e
	}

	r.content = b
	return b, nil
}

// 获取JSON响应内容
// 可以传指针类型的接收者
func (r *GapiResponse) JSON(v ...interface{}) (interface{}, error) {
	b, err := r.Content()
	if err != nil {
		return nil, err
	}

	if !strings.Contains(r.Header.Get(HeaderContentType), "json") {
		err := r.Status
		if len(b) > 0 {
			err = string(b)
		}
		return nil, errors.New(err)
	}

	var res interface{}
	if len(v) > 0 {
		res = v[0]
	} else {
		res = new(map[string]interface{})
	}

	if err = json.Unmarshal(b, res); err != nil {
		return nil, err
	}

	return res, nil
}

// 获取JSON响应内容
// 可以传指针类型的接收者
func (r *GapiResponse) GJSON() (*gjson.Result, error) {
	if r.result != nil {
		return r.result, nil
	}

	b, err := r.Content()
	if err != nil {
		return nil, err
	}

	if !strings.Contains(r.Header.Get(HeaderContentType), "json") {
		err := r.Status
		if len(b) > 0 {
			err = string(b)
		}
		return nil, errors.New(err)
	}

	res := gjson.ParseBytes(b)
	r.result = &res

	return r.result, nil
}

// 获取文字响应内容
func (r *GapiResponse) Text() (string, error) {
	b, err := r.Content()

	if err != nil {
		return "", err
	}

	return string(b), nil
}

// 获取最终请求的URL
func (r *GapiResponse) URL() (*url.URL, error) {
	u := r.Request.URL

	if r.StatusCode == http.StatusMovedPermanently ||
		r.StatusCode == http.StatusFound ||
		r.StatusCode == http.StatusSeeOther ||
		r.StatusCode == http.StatusTemporaryRedirect {
		location, err := r.Location()

		if err != nil {
			return nil, err
		}

		u = u.ResolveReference(location)
	}

	return u, nil
}

// 获取相应代码的描述
func (r *GapiResponse) Reason() string {
	return http.StatusText(r.StatusCode)
}

// 判断响应是否成功
// 其实就是判断响应状态码是否在100~399之间
func (r *GapiResponse) OK() bool {
	return r.StatusCode < 400
}

func (r *GapiResponse) Gapi() *Gapi {
	return r.p
}
