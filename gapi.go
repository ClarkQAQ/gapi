package gapi

import (
	"net/http"
	"net/url"
	"sync"
	"time"

	"utilware/cache"
)

var (
	defaultGapiURL, _ = url.Parse("https://www.pixiv.net") // 默认的pixiv.net的URL
	defaultTimeout    = time.Second * 15                   // 默认的超时时间
)

// Gapi 客户端
type Gapi struct {
	c           *http.Client
	siteURL     *url.URL
	timeout     time.Duration
	gHeaderLock *sync.RWMutex
	gHeader     http.Header
	cache       *cache.Cache
}

// 创建客户端选项
type Options struct {
	ProxyURL string
	Timeout  time.Duration
}

// New 创建一个Gapi客户端
// 可以通过设置选项来设置客户端的一些选项
func New(siteURL string, opts *Options) (*Gapi, error) {
	u, e := url.Parse(siteURL)
	if e != nil {
		return nil, e
	}

	p := &Gapi{
		c:           &http.Client{},
		siteURL:     u,
		timeout:     defaultTimeout,
		gHeaderLock: &sync.RWMutex{},
		gHeader:     http.Header{},
		cache:       cache.New(5*time.Minute, 10*time.Minute),
	}

	p.ClearCookies()
	p.SetOptions(opts)

	return p, nil
}

func (p *Gapi) SetOptions(opts *Options) error {
	if opts == nil {
		return nil
	}

	if opts.ProxyURL != "" {
		p.newProxy(opts.ProxyURL)
	}

	return nil
}

func (p *Gapi) Timeout() time.Duration {
	return p.timeout
}

func (p *Gapi) Client() *http.Client {
	return p.c
}

func (p *Gapi) GetURL() *url.URL {
	return p.siteURL
}

func (p *Gapi) Cache() *cache.Cache {
	return p.cache
}

func (p *Gapi) SetGHeader(m http.Header) {
	p.gHeaderLock.Lock()
	defer p.gHeaderLock.Unlock()

	p.gHeader = m
}

func (p *Gapi) GetGHeader() http.Header {
	p.gHeaderLock.RLock()
	defer p.gHeaderLock.RUnlock()

	return p.gHeader
}
