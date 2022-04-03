package gpixiv

import (
	"net/http"
	"net/url"
	"time"
)

var (
	defaultPixivURL, _ = url.Parse("https://www.pixiv.net")
	defaultTimeout     = time.Second * 10
	defaultUserAgent   = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36"
)

type Pixiv struct {
	c         *http.Client
	pixivURL  *url.URL
	timeout   time.Duration
	userAgent string
}

type Options struct {
	URL       string
	ProxyURL  string
	UserAgent string
	Timeout   time.Duration
}

func NewClient(opts *Options) (*Pixiv, error) {
	p := &Pixiv{
		c:         &http.Client{},
		pixivURL:  defaultPixivURL,
		timeout:   defaultTimeout,
		userAgent: defaultUserAgent,
	}

	if opts == nil {
		return p, nil
	}

	if opts.URL != "" {
		u, e := url.Parse(opts.URL)
		if e != nil {
			return nil, e
		}

		p.pixivURL = u
	}

	if opts.ProxyURL != "" {
		p.newProxy(opts.ProxyURL)
	}

	if opts.Timeout > 0 {
		p.timeout = opts.Timeout
	}

	if opts.UserAgent != "" {
		p.userAgent = opts.UserAgent
	}

	return p, nil
}
