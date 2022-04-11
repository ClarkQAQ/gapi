package pixiv

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
	"utilware/goquery"

	"github.com/ClarkQAQ/gapi"
)

// 获取收藏artwork/获取bookmark
// tag => 插画/漫画收藏标签
// rest => 	show:公开模式/hide:非公开模式
// lang => 返回语言
// offset => 偏移数量
// limit => 每页数量
func BookmarkList(tag string, rest string, lang string, offset int64, limit int64) *gapi.GapiApi {
	a := gapi.NewAPI("GET", "").
		SetValues(map[string]string{
			"tag":    tag,
			"offset": fmt.Sprintf("%d", offset),
			"limit":  fmt.Sprintf("%d", limit),
			"rest":   rest,
			"lang":   lang,
		}).
		SetHeader(gapi.HeaderAccept, "application/json")

	a.SetHijack(func(p *gapi.Gapi, req *http.Request) error {
		userId, ok := p.Cache().Get("user_id")
		if !ok || userId == nil {
			resp, e := p.Do(GetUserInfo())
			if e != nil {
				return fmt.Errorf("get user info: %s", e.Error())
			}

			userInfo := GetUserInfoData{}
			if _, e := resp.JSON(&userInfo); e != nil {
				return fmt.Errorf("parse user info: %s", e.Error())
			}

			userId = userInfo.Id

			p.Cache().Set("user_id", userId, time.Minute*30)
		}

		req.URL.Path = fmt.Sprintf("/ajax/user/%v/illusts/bookmarks", userId)
		return nil
	})

	a.SetRespHijack(func(resp *gapi.GapiResponse, setBody func(body []byte)) error {
		g, e := resp.GJSON()
		if e != nil {
			return fmt.Errorf("parse json: %s", e.Error())
		}

		if !g.Get("error").Exists() {
			return errors.New("response content error")
		}

		if g.Get("error").Bool() {
			return fmt.Errorf("response error: %s", g.Get("message").String())
		}

		setBody([]byte(g.Get("body").Raw))
		return nil
	})

	return a
}

// 收藏artwork/添加bookmark
// illust_id => 插画/漫画编号
// restrict => 0:公开模式/1:非公开模式
// comment => 收藏备注
// tags => 收藏标签
func AddBookmark(illust_id int64, restrict int, comment string, tags []string) *gapi.GapiApi {
	// 定义基础
	a := gapi.NewAPI("POST", "/ajax/illusts/bookmarks/add").
		SetHeader(gapi.HeaderAccept, "application/json; charset=utf-8").
		SetHeader(gapi.HeaderContentType, "application/json; charset=utf-8")

	// 定义内容
	m := make(map[string]interface{}, 4)
	if illust_id <= 0 {
		a.SetError(errors.New("illust_id is invalid"))
		return a
	}

	a.SetHijack(func(p *gapi.Gapi, req *http.Request) error {
		// 这个接口要求csrf token
		csrfToken, e := GetCsrfTokenString(p)
		if e != nil {
			return e
		}

		req.Header.Set(gapi.HeaderXCSRFToken, csrfToken)
		return nil
	})

	m["illust_id"] = illust_id
	m["restrict"] = restrict
	m["comment"] = comment
	m["tags"] = tags

	// 序列化为json
	b, e := json.Marshal(m)
	if e != nil {
		a.SetError(e)
		return a
	}

	a.SetBody(b)
	a.SetRespHijack(func(resp *gapi.GapiResponse, setBody func(body []byte)) error {
		res, e := resp.GJSON()
		if e != nil {
			return e
		}

		if !res.Get("error").Exists() {
			return errors.New("response error")
		}

		if res.Get("error").Bool() {
			return errors.New(res.Get("message").String())
		}

		setBody([]byte("success"))
		return nil
	})

	return a
}

// 取消收藏artwork/取消添加bookmark
// illust_id => 插画/漫画编号
func DeleteBookmark(illust_id int64) *gapi.GapiApi {
	// 访问bookmark页面获取bookmarkId
	a := gapi.NewAPI("GET", "/bookmark_add.php").
		SetValue("type", "illust").
		SetValue("illust_id", fmt.Sprint(illust_id)).
		SetHeader(gapi.HeaderAccept, "text/html")

	a.SetRespHijack(func(resp *gapi.GapiResponse, setBody func(body []byte)) error {
		if resp.StatusCode != 200 {
			return fmt.Errorf("status code is %d", resp.StatusCode)
		}

		// 加载进Goquery, 然后解析 bookmark 页面
		doc, e := goquery.NewDocumentFromReader(bytes.NewReader(resp.Raw()))
		if e != nil {
			return fmt.Errorf("goquery load document: %s", e.Error())
		}

		resp.Header.Set(gapi.HeaderContentType, "application/json")

		// 寻找 bookmark id, 如果找不到多半是因为没有收藏
		bookmarkId, ok := doc.Find(".remove-bookmark-form > input:nth-child(5)").Attr("value")
		if !ok {
			setBody([]byte(`{"error":true,"message":"illust is not favorited","body":[]}`))
			return nil
		}

		// 调用接口取消收藏
		a := gapi.NewAPI("POST", "/rpc/index.php").
			SetHeader(gapi.HeaderAccept, "application/json; charset=utf-8").
			// 这个接口需要特殊请求数据
			SetHeader(gapi.HeaderContentType, "application/x-www-form-urlencoded; charset=utf-8").
			SetHijack(func(p *gapi.Gapi, req *http.Request) error {
				// 也需要csrf token
				csrfToken, e := GetCsrfTokenString(p)
				if e != nil {
					return e
				}

				req.Header.Set(gapi.HeaderXCSRFToken, csrfToken)

				return nil
			})

		// 设置请求内容
		a.SetBody([]byte(url.Values{
			"mode":        {"delete_illust_bookmark"},
			"bookmark_id": {bookmarkId},
		}.Encode()))

		resp2, e := resp.Gapi().Do(a)
		if e != nil {
			return e
		}

		res, e := resp2.GJSON()
		if e != nil {
			return e
		}

		if !res.Get("error").Exists() {
			return errors.New("response error")
		}

		if res.Get("error").Bool() {
			return errors.New(res.Get("message").String())
		}

		setBody([]byte("success"))
		return nil
	})

	return a
}
