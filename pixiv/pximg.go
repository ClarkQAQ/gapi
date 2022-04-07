package pixiv

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/ClarkQAQ/gapi"
)

// 获取图片数据
// picURL => 图片URL
func Pximg(picURL string) *gapi.GapiApi {
	a := gapi.NewAPI("GET", picURL)

	u, e := url.Parse(picURL)
	if e != nil {
		a.SetError(e)
		return a
	}

	a.SetHeader(gapi.HeaderReferer, fmt.Sprintf("%s://%s", u.Scheme, u.Host)).
		SetHeader(gapi.HeaderAccept, "image/webp,image/apng,image/*,*/*;q=0.8").
		SetHeader("Host", u.Host).
		SetHeader("Upgrade-Insecure-Requests", "1").
		SetRespHijack(func(resp *gapi.GapiResponse, setBody func([]byte)) error {
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("server response status code: %d", resp.StatusCode)
			}

			return nil
		})

	return a
}
