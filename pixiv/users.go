package pixiv

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"utilware/gjson"
	"utilware/goquery"

	"github.com/ClarkQAQ/gapi"
)

// 使用Cookie登录
// phpsessid => PHPSESSID (浏览器f12获取)
func CookieLogin(phpsessid string) *gapi.GapiApi {
	a := gapi.NewAPI("HEAD", "/touch/ajax/follow/latest").
		SetHeader("Accept", "application/json; charset=utf-8").
		SetHijack(func(p *gapi.Gapi, req *http.Request) error {
			p.ClearCookies()
			p.Client().Jar.SetCookies(req.URL, []*http.Cookie{{
				Domain: "." + p.GetURL().Host,
				Path:   "/",
				Name:   "PHPSESSID",
				Value:  phpsessid,
			}})

			return CheckPHPSESSID(p)
		}).SetRespHijack(func(resp *gapi.GapiResponse, setBody func(body []byte)) error {
		if resp.StatusCode > 200 {
			return fmt.Errorf("login failed: %d", resp.StatusCode)
		}

		setBody([]byte("success"))
		return nil
	})

	return a
}

type GetUserInfoData struct {
	Id            string `json:"id"`
	GapiId        string `json:"pixiv_id"`
	Name          string `json:"name"`
	ProfileImg    string `json:"profile_img"`
	ProfileImgBig string `json:"profile_img_big"`
	IsPremium     bool   `json:"is_premium"`
	XRestrict     int    `json:"xrestrict"`
	Adult         bool   `json:"adult"`
	SafeMode      bool   `json:"safe_mode"`
	Development   bool   `json:"development"`
	CsrfToken     string `json:"csrf_token"`
}

// 获取账户的设置
func GetUserInfo() *gapi.GapiApi {
	return gapi.NewAPI("GET", "/setting_user.php").
		SetHijack(func(p *gapi.Gapi, req *http.Request) error {
			return CheckPHPSESSID(p)
		}).
		SetHeader("Accept", "*/*").
		SetRespHijack(func(resp *gapi.GapiResponse, setBody func([]byte)) error {
			// 加载进Goquery
			doc, e := goquery.NewDocumentFromReader(bytes.NewReader(resp.Raw()))
			if e != nil {
				return fmt.Errorf("goquery load document: %s", e.Error())
			}

			// 获取input content
			metaData, ok := doc.Find("#meta-global-data").Attr("content")
			if !ok {
				return fmt.Errorf("meta-global-data not found")
			}

			// 设置响应头
			resp.Response.Header.Set(gapi.HeaderContentType, "application/json")

			gjson.Get(metaData, "userData")

			userData := gjson.Get(metaData, "userData")
			m := GetUserInfoData{
				Id:            userData.Get("id").String(),
				GapiId:        userData.Get("pixivId").String(),
				Name:          userData.Get("name").String(),
				ProfileImg:    userData.Get("profileImg").String(),
				ProfileImgBig: userData.Get("profileImgBig").String(),
				IsPremium:     userData.Get("premium").Bool(),
				XRestrict:     int(userData.Get("xRestrict").Int()),
				Adult:         userData.Get("adult").Bool(),
				SafeMode:      userData.Get("safeMode").Bool(),
				CsrfToken:     gjson.Get(metaData, "token").String(),
			}

			b, e := json.Marshal(m)
			if e != nil {
				return e
			}

			// 设置响应体
			setBody(b)
			return nil
		})
}

// 获取Csrf Token
func GetCsrfToken() *gapi.GapiApi {
	return gapi.NewAPI("GET", "/setting_user.php").
		SetHijack(func(p *gapi.Gapi, req *http.Request) error {
			return CheckPHPSESSID(p)
		}).
		SetHeader("Accept", "*/*").
		SetRespHijack(func(resp *gapi.GapiResponse, setBody func([]byte)) error {
			// 加载进Goquery
			doc, e := goquery.NewDocumentFromReader(bytes.NewReader(resp.Raw()))
			if e != nil {
				return fmt.Errorf("goquery load document: %s", e.Error())
			}

			// 获取input content
			metaData, ok := doc.Find("#meta-global-data").Attr("content")
			if !ok {
				return fmt.Errorf("meta-global-data not found")
			}

			// 设置响应头
			resp.Response.Header.Set(gapi.HeaderContentType, "text/html; charset=utf-8")

			// 获取csrf_token
			csrfToken := gjson.Get(metaData, "token").String()

			// 设置响应体
			setBody([]byte(csrfToken))

			return nil
		})
}

// 此处代码已废弃...
// 妈的被白嫖了...这种登录方式特别难
// 要么模拟浏览器,要么只能手动设置PHPSESSID
// func (p *Gapi) Login(username, password string) (string, error) {
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

// func (p *Gapi) getLoginPage() ([]byte, error) {
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

// func (p *Gapi) getLoginPostKey() (string, error) {
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
