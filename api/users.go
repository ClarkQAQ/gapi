package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gpixiv"
	"net/http"
	"utilware/gjson"
	"utilware/goquery"
)

type GetUserInfoData struct {
	Id            string `json:"id"`
	PixivId       string `json:"pixiv_id"`
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
func GetUserInfo() *gpixiv.PixivApi {
	return gpixiv.NewAPI("GET", "/setting_user.php").
		SetHijack(func(p *gpixiv.Pixiv, req *http.Request) error {
			return CheckPHPSESSID(p)
		}).
		SetHeader("Accept", "*/*").
		SetRespHijack(func(resp *gpixiv.PixivResponse, setBody func([]byte)) error {
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
			resp.Response.Header.Set(gpixiv.HeaderContentType, "application/json")

			gjson.Get(metaData, "userData")

			userData := gjson.Get(metaData, "userData")
			m := GetUserInfoData{
				Id:            userData.Get("id").String(),
				PixivId:       userData.Get("pixivId").String(),
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
func GetCsrfToken() *gpixiv.PixivApi {
	return gpixiv.NewAPI("GET", "/setting_user.php").
		SetHijack(func(p *gpixiv.Pixiv, req *http.Request) error {
			return CheckPHPSESSID(p)
		}).
		SetHeader("Accept", "*/*").
		SetRespHijack(func(resp *gpixiv.PixivResponse, setBody func([]byte)) error {
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
			resp.Response.Header.Set(gpixiv.HeaderContentType, "text/html; charset=utf-8")

			// 获取csrf_token
			csrfToken := gjson.Get(metaData, "token").String()

			// 设置响应体
			setBody([]byte(csrfToken))

			return nil
		})
}
