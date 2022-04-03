package api

import (
	"net/http"
	"net/url"
)

type PixivApi struct {
	c        *http.Client
	pixivURL *url.URL
}
