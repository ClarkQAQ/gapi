package ccoo

import (
	"net/http"

	"github.com/ClarkQAQ/gapi"
)

const (
	URL = ""
)

var (
	GlobalHeader = http.Header{
		"User-Agent": {gapi.DefaultUserAgent},
	}
)
