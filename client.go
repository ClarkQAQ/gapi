package gpixiv

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

// 生成URL
func (p *Pixiv) EndpointURL(url string, values url.Values) *url.URL {
	u := *p.pixivURL
	u.Path = url

	if values != nil {
		u.RawQuery = values.Encode()
	}

	return &u
}

func (p *Pixiv) Request(ctx context.Context, method, url string, ref func(req *http.Request) (*http.Request, error)) (*http.Response, error) {
	req, e := http.NewRequestWithContext(ctx, "GET", url, nil)
	if e != nil {
		return nil, e
	}

	if ref != nil {
		req, e = ref(req)
		if e != nil {
			return nil, e
		}
	}

	return p.c.Do(req)
}

func (p *Pixiv) GetPximg(picURL string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	u, e := url.Parse(picURL)
	if e != nil {
		return nil, e
	}

	req, e := p.Request(ctx, "GET", picURL, func(req *http.Request) (*http.Request, error) {
		req.Header.Set("Referer", fmt.Sprintf("%s://%s",
			u.Scheme, u.Host))
		req.Header.Set("User-Agent", p.userAgent)
		req.Header.Set("Accept", "image/webp,image/apng,image/*,*/*;q=0.8")
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("Host", u.Host)
		req.Header.Set("Upgrade-Insecure-Requests", "1")

		return req, nil
	})

	if e != nil {
		return nil, e
	}

	if req.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server response status code: %d", req.StatusCode)
	}

	b, e := ioutil.ReadAll(req.Body)
	if e != nil {
		return nil, e
	}

	return b, nil
}

// 清除 cookie
func (p *Pixiv) clearCookies() {
	p.c.Jar, _ = cookiejar.New(nil)
}
