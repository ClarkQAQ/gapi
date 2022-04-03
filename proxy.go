package gpixiv

import (
	"net/http"
	"net/url"

	"golang.org/x/net/proxy"
)

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
