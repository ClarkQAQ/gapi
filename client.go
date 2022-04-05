package gpixiv

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

// 生成URL
// 提供Pixiv站点的path和参数返回完整的URL
func (p *Pixiv) EndpointPATH(path string, values url.Values) *url.URL {
	u := *p.pixivURL
	u.Path = path

	if values != nil {
		u.RawQuery = values.Encode()
	}

	return &u
}

// 生成URL
// 提供任意的url和参数返回完整的URL
func (p *Pixiv) EndpointURL(urlString string, values url.Values) (*url.URL, error) {
	u, e := url.Parse(urlString)
	if e != nil {
		return nil, e
	}

	if values != nil {
		u.RawQuery = values.Encode()
	}

	return u, nil
}

// 创建一个新的http.Client
// 并且有一个ref可以设置自定义选项然后通过client.Do()来发送请求
// 将返回原始的http.Response
func (p *Pixiv) Request(ctx context.Context, method, url string, body io.Reader, hijack func(c *http.Client, req *http.Request) error) (*http.Response, error) {
	req, e := http.NewRequestWithContext(ctx, method, url, body)
	if e != nil {
		return nil, e
	}

	// cp := *p.c
	// c := &cp
	c := p.c

	req.Header.Set(HeaderUserAgent, p.userAgent)

	if hijack != nil {
		if e = hijack(c, req); e != nil {
			return nil, e
		}
	}

	defer c.CloseIdleConnections()

	return c.Do(req)
}

// 获取图片比特数据
// 传入图片的URL, 返回图片的比特数据
// 理论上Pixiv大部分图片都是支持的
func (p *Pixiv) Pximg(picURL string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	u, e := url.Parse(picURL)
	if e != nil {
		return nil, e
	}

	resp, e := p.Request(ctx, "GET", picURL, nil, func(c *http.Client, req *http.Request) error {
		req.Header.Set("Referer", fmt.Sprintf("%s://%s", u.Scheme, u.Host))
		req.Header.Set("User-Agent", p.userAgent)
		req.Header.Set("Accept", "image/webp,image/apng,image/*,*/*;q=0.8")
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("Host", u.Host)
		req.Header.Set("Upgrade-Insecure-Requests", "1")

		return nil
	})

	if e != nil {
		return nil, e
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server response status code: %d", resp.StatusCode)
	}

	b, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return nil, e
	}

	return b, nil
}

// 清除 cookie
// 将cookiejar中的所有cookie清除
// 当设置了PHPSESSID时, 可以使用这个方法来清除
func (p *Pixiv) clearCookies() {
	p.c.Jar, _ = cookiejar.New(nil)
}
