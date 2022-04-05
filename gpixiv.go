package gpixiv

import (
	"net/http"
	"net/url"
	"time"
)

var (
	defaultPixivURL, _ = url.Parse("https://www.pixiv.net") // 默认的pixiv.net的URL
	defaultTimeout     = time.Second * 15                   // 默认的超时时间
)

// Pixiv 客户端
type Pixiv struct {
	c         *http.Client
	pixivURL  *url.URL
	timeout   time.Duration
	userAgent string
	language  string
}

// 创建客户端选项
type Options struct {
	URL       string
	ProxyURL  string
	UserAgent string
	Language  string
	Timeout   time.Duration
}

// New 创建一个Pixiv客户端
// 可以通过设置选项来设置客户端的一些选项
func New(opts *Options) (*Pixiv, error) {
	p := &Pixiv{
		c:         &http.Client{},
		pixivURL:  defaultPixivURL,
		timeout:   defaultTimeout,
		userAgent: defaultUserAgent,
		language:  defaultLanguage,
	}

	p.SetOptions(opts)

	return p, nil
}

func (p *Pixiv) SetOptions(opts *Options) error {
	if opts == nil {
		return nil
	}

	if opts.URL != "" {
		u, e := url.Parse(opts.URL)
		if e != nil {
			return e
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

	if opts.Language != "" {
		p.language = opts.Language
	}

	return nil
}
