package gpixiv

import (
	"context"
	"net/http"
)

// 判断是否登录 (需要PHPSESSID)
// 原理是访问用户设置页面, 如果没有302到登录页面, 则已登录
// 当然还有一个判断是否404的分支
func (p *Pixiv) IsLogged() (ret bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	resp, e := p.Request(ctx, http.MethodHead, p.EndpointPATH("/touch/ajax/follow/latest", nil).String(), nil, func(c *http.Client, req *http.Request) error {
		c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
		return nil
	})

	if e != nil {
		return false, e
	}
	defer resp.Body.Close()

	return resp.StatusCode <= http.StatusOK, nil
}

// 设置 PHPSESSID
// 暂时没有自动登录的功能, 只能手动设置
func (p *Pixiv) SetPHPSESSID(sessionID string) {
	p.clearCookies()

	p.c.Jar.SetCookies(
		p.pixivURL,
		[]*http.Cookie{{
			Domain: "." + p.pixivURL.Host,
			Path:   "/",
			Name:   "PHPSESSID",
			Value:  sessionID,
		}},
	)
}

// 此处代码已废弃...
// 妈的被白嫖了...这种登录方式特别难
// 要么模拟浏览器,要么只能手动设置PHPSESSID
// func (p *Pixiv) Login(username, password string) (string, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
// 	defer cancel()

// 	u, e := p.EndpointURL("https://accounts.pixiv.net/api/login", url.Values{
// 		"lang": {p.language},
// 	})
// 	if e != nil {
// 		return "", e
// 	}

// 	postKey, e := p.getLoginPostKey()
// 	if e != nil {
// 		return "", e
// 	}

// 	fromData := url.Values{}
// 	fromData.Set("pixiv_id", username)
// 	fromData.Set("password", password)
// 	fromData.Set("captcha", "")
// 	fromData.Set("g_recaptcha_response", "")
// 	fromData.Set("post_key", postKey)
// 	fromData.Set("return_to", "https://www.pixiv.net/")
// 	fromData.Set("source", "accounts")
// 	fromData.Set("app_ios", "0")
// 	fromData.Set("ref", "")

// 	fmt.Println(fromData.Encode())

// 	resp, e := p.Request(ctx, http.MethodPost, u.String(), strings.NewReader(fromData.Encode()),
// 		func(c *http.Client, req *http.Request) error {
// 			req.Header.Del(HeaderUserAgent)
// 			req.Header.Set(HeaderUserAgent, IosDeviceUserAgent)
// 			req.Header.Set(HeaderContentType, HeaderXFronURLEncoded)

// 			return nil
// 		})
// 	if e != nil {
// 		return "", e
// 	}

// 	b, e := ioutil.ReadAll(resp.Body)
// 	if e != nil {
// 		return "", e
// 	}

// 	fmt.Println(string(b))
// 	g := gjson.ParseBytes(b)

// 	fmt.Println(g.Get("body.validation_errors.etc").String())
// 	if g.Get("error").Bool() {
// 		return "", errors.New(g.Get("message").String())
// 	}

// 	return "", nil
// }

// func (p *Pixiv) getLoginPage() ([]byte, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
// 	defer cancel()

// 	u, e := p.EndpointURL("https://accounts.pixiv.net/login", url.Values{
// 		"source":    {"touch"},
// 		"return_to": {"https://www.pixiv.net/"},
// 		"view_type": {"page"},
// 		"lang":      {p.language},
// 	})
// 	if e != nil {
// 		return nil, e
// 	}

// 	resp, err := p.Request(ctx, http.MethodGet, u.String(), nil,
// 		func(c *http.Client, req *http.Request) error {
// 			req.Header.Del(HeaderUserAgent)
// 			req.Header.Set(HeaderUserAgent, IosDeviceUserAgent)

// 			return nil
// 		})
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	return ioutil.ReadAll(resp.Body)
// }

// func (p *Pixiv) getLoginPostKey() (string, error) {
// 	pageBytes, e := p.getLoginPage()
// 	if e != nil {
// 		return "", e
// 	}

// 	doc, e := goquery.NewDocumentFromReader(bytes.NewReader(pageBytes))
// 	if e != nil {
// 		return "", e
// 	}

// 	config := doc.Find("#init-config").AttrOr("value", "")
// 	if config == "" {
// 		return "", errors.New("pixiv: client: login: init-config not found")
// 	}

// 	postKey := gjson.Get(config, `pixivAccount\.postKey`).String()
// 	if postKey == "" {
// 		return "", errors.New("pixiv: client: login: postKey not found")
// 	}

// 	return postKey, nil
// }
