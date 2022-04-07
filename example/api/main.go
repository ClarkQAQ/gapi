package main

import (
	"fmt"
	"utilware/logger"

	"github.com/ClarkQAQ/gapi"
	"github.com/ClarkQAQ/gapi/pixiv"
)

func main() {
	p, e := gapi.New(pixiv.URL, &gapi.Options{
		ProxyURL: "socks5://127.0.0.1:7891",
	})
	if e != nil {
		logger.Fatal("创建Gapi客户端失败: %s", e.Error())
	}

	resp, e := p.Do(CustomApi(1))
	if e != nil {
		logger.Fatal("调用Gapi客户端失败: %s", e.Error())
	}

	logger.Info("Hijack结果: %s", resp.Raw())
}

func CustomApi(id int64) *gapi.GapiApi {
	a := gapi.NewAPI("GET", fmt.Sprintf("/ajax/illust/%d/pages", id)).
		SetHeader("Accept", "application/json; charset=utf-8").
		SetValue("lang", "en").
		// RespHijack 拦截响应
		SetRespHijack(func(resp *gapi.GapiResponse, setBody func(body []byte)) error {
			logger.Info("原始响应状态: %s", resp.Status)

			setBody([]byte("Hijacked"))
			return nil
		})

	return a
}
