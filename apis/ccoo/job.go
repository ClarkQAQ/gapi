package ccoo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"utilware/goquery"

	"github.com/ClarkQAQ/gapi"
)

type JobListItem struct {
	Title       string   `json:"title"`        // 标题
	Link        string   `json:"link"`         // 链接
	Price       string   `json:"price"`        // 薪资
	WelfareTags []string `json:"welfare_tags"` // 福利标签???
	CompanyName string   `json:"company_name"` // 公司名称
	CompanyAddr string   `json:"company_addr"` // 公司地址
	PopRequire  []string `json:"pop_require"`  // 公司要求
	UpdateDate  string   `json:"update_date"`  // 更新日期
}

// 获取招聘信息列表
// curl => 网站内容链接
// 一般是网页的: http://www.heshan.ccoo.cn/post/zhaopins/list-0-0-0-0-0-0-%d.html
// %d => 页码 (从0开始) 当然其他地方不确定可能要自己fork一下了
// 这两段 `heshan.ccoo.cn` `0-0-0-0-0-0-1` 具体含义貌似每个地区都不一样, 所以只能先这样抽了
func Jobs(curl string) *gapi.GapiApi {
	// 解析链接
	u, e := url.Parse(curl)

	// 创建一个用于返回的 GapiApi
	a := gapi.NewAPI("GET", curl).
		SetHeader(gapi.HeaderAccept, "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/jxl,image/webp,*/*;q=0.8").
		SetHeader(gapi.HeaderReferer, curl)

	// 链接错误返回
	if e != nil {
		a.SetError(fmt.Errorf("parse curl error: %s", e.Error()))
		return a
	}

	// 内容Hijack
	a.SetRespHijack(func(resp *gapi.GapiResponse, setBody func(body []byte)) error {
		// 加载进Goquery, 然后解析 bookmark 页面
		doc, e := goquery.NewDocumentFromReader(bytes.NewReader(resp.Raw()))
		if e != nil {
			return fmt.Errorf("goquery load document: %s", e.Error())
		}

		// 找到列表
		n := doc.Find("li.item").Nodes
		if len(n) <= 0 {
			return errors.New("goquery load document: not found job list")
		}

		retData := make([]JobListItem, len(n))
		nowDate := time.Now().Format("2006-01-02")

		// 曲线救国, 把每一个列表信息都解析出来
		splitList := func(s string) []string {
			list := strings.Split(strings.ReplaceAll(s, " ", ""), "\n")
			for i := 0; i < len(list); i++ {
				if list[i] == "" {
					list = append(list[:i], list[i+1:]...)
					i--
				}
			}

			return list
		}

		for i := 0; i < len(n); i++ {
			node := goquery.NewDocumentFromNode(n[i])
			link, _ := node.Find(".jobInfo > .clearfix > a").Attr("href")

			ret := JobListItem{
				Title:       node.Find(".jobInfo > .clearfix > a > p").Text(),
				Link:        fmt.Sprintf("%s://%s", u.Scheme, path.Join(u.Host, link)),
				Price:       node.Find(".jobInfo > .price").Text(),
				WelfareTags: splitList(node.Find(".jobInfo > .listBox").Text()),
				CompanyName: node.Find(".comInfo > .clearfix > .tit").Text(),
				CompanyAddr: node.Find(".comInfo > .addr").Text(),
				PopRequire:  splitList(node.Find(".comInfo > .listBox").Text()),
				UpdateDate:  node.Find(".dateBox > .date").Text(),
			}

			// 更新日期是今天的就转换成日期
			if ret.UpdateDate == "今天" {
				ret.UpdateDate = nowDate
			}

			retData[i] = ret
		}

		b, e := json.Marshal(retData)
		if e != nil {
			return fmt.Errorf("json marshal: %s", e.Error())
		}

		// 设置响应头
		resp.Response.Header.Set(gapi.HeaderContentType, "application/json")

		setBody(b)
		return nil
	})

	return a
}

// 获取招聘信息列表
// curl => 网站API链接
// 我本地的是: http://www.heshan.ccoo.cn/home/recruit/getpositiontel.html
// id => Job 的 ID, 链接后面的数字 (好像要去掉x)
func Phone(curl string, id string) *gapi.GapiApi {
	// 解析链接
	u, e := url.Parse(curl)

	// 创建一个用于返回的 GapiApi
	a := gapi.NewAPI("POST", curl).
		SetHeader(gapi.HeaderAccept, "application/json, text/javascript, */*; q=0.01").
		SetHeader(gapi.HeaderContentType, "application/x-www-form-urlencoded; charset=UTF-8").
		SetHeader("X-Requested-With", "XMLHttpRequest").
		SetHeader(gapi.HeaderReferer, curl)

	// 链接错误返回
	if e != nil {
		a.SetError(fmt.Errorf("parse curl error: %s", e.Error()))
		return a
	}

	a.SetHeader(gapi.HeaderOrigin, u.Scheme+"://"+u.Host).
		SetHijack(func(p *gapi.Gapi, req *http.Request) error {
			if v, b := p.Cache().Get("cookie"); b {
				req.Header.Set("Cookie", fmt.Sprint(v))
			}

			return nil
		})

	b := url.Values{
		"id":    {id},
		"otype": {"0"},
	}.Encode()

	a.SetBody([]byte(b))

	// 内容Hijack
	a.SetRespHijack(func(resp *gapi.GapiResponse, setBody func(body []byte)) error {
		res, e := resp.GJSON()
		if e != nil {
			return fmt.Errorf("parse json: %s", e.Error())
		}

		phone := res.Get("ServerInfo").String()
		if phone == "" {
			return errors.New("not found phone")
		}

		setBody([]byte(phone))
		return nil
	})

	return a
}
