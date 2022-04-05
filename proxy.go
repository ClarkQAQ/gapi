package gpixiv

import (
	"net/http"
	"net/url"

	"utilware/dep/x/net/proxy"
)

// 国区特供代理功能
// 可以通过设置代理的URL来在国内访问Pixiv站点
// 现在支持的代理有: socks5, http
func (p *Pixiv) newProxy(proxyURL string) error {
	u, e := url.Parse(proxyURL)
	if e != nil {
		return e
	}

	switch u.Scheme {
	case "socks5":
		p.newSocks5Proxy(u)
	case "http":
		p.newHttpProxy(u)
	}

	return nil
}

// 新建socks5代理并设置到http客户端
func (p *Pixiv) newSocks5Proxy(u *url.URL) error {
	var auth *proxy.Auth = nil
	if u.User != nil {
		password, _ := u.User.Password()
		auth = &proxy.Auth{User: u.User.Username(), Password: password}
	}

	dialer, e := proxy.SOCKS5("tcp", u.Host, auth, proxy.Direct)
	if e != nil {
		return e
	}

	httpTransport := &http.Transport{}
	p.c.Transport = httpTransport
	httpTransport.Dial = dialer.Dial
	return nil
}

// 新建一个http代理并设置到http客户端
func (p *Pixiv) newHttpProxy(u *url.URL) error {
	dialer, e := proxy.FromURL(u, proxy.Direct)
	if e != nil {
		return e
	}

	httpTransport := &http.Transport{}
	p.c.Transport = httpTransport
	httpTransport.Dial = dialer.Dial
	return nil
}
