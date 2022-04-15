package gapi

import (
	"context"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

// 生成URL
// 提供Gapi站点的path和参数返回完整的URL
func (p *Gapi) EndpointPATH(path string, values url.Values) *url.URL {
	u := *p.siteURL
	u.Path, u.RawQuery, _ = strings.Cut(path, "?")

	if values != nil && len(values) > 0 {
		u.RawQuery = values.Encode()
	}

	return &u
}

// 生成URL
// 提供任意的url和参数返回完整的URL
func (p *Gapi) EndpointURL(urlString string, values url.Values) (*url.URL, error) {
	u, e := url.Parse(urlString)
	if e != nil {
		return nil, e
	}

	if values != nil && len(values) > 0 {
		u.RawQuery = values.Encode()
	}

	return u, nil
}

// 创建一个新的http.Client
// 并且有一个ref可以设置自定义选项然后通过client.Do()来发送请求
// 将返回原始的http.Response
func (p *Gapi) Request(ctx context.Context, method, url string, body io.Reader, hijack func(c *http.Client, req *http.Request) error) (*http.Response, error) {
	req, e := http.NewRequestWithContext(ctx, method, url, body)
	if e != nil {
		return nil, e
	}

	p.gHeaderLock.RLock()
	for k, v := range p.gHeader {
		req.Header[k] = v
	}
	p.gHeaderLock.RUnlock()

	// logger.Info("%+v", req)

	if hijack != nil {
		if e = hijack(p.c, req); e != nil {
			return nil, e
		}
	}

	defer p.c.CloseIdleConnections()

	return p.c.Do(req)
}

// 清除 cookie
// 将cookiejar中的所有cookie清除
// 当设置了PHPSESSID时, 可以使用这个方法来清除
func (p *Gapi) ClearCookies() {
	p.c.Jar, _ = cookiejar.New(nil)
}
