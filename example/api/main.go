package main

import (
	"fmt"
	"gpixiv"
	"utilware/logger"
)

func main() {
	p, e := gpixiv.New(&gpixiv.Options{
		ProxyURL: "socks5://127.0.0.1:7891",
	})
	if e != nil {
		logger.Fatal("创建Pixiv客户端失败: %s", e.Error())
	}

	resp, e := p.Do(CustomApi(1))
	if e != nil {
		logger.Fatal("调用Pixiv客户端失败: %s", e.Error())
	}

	// 输出结果
	// 无论如何结果都是Hijacked
	text, e := resp.Text()
	if e != nil {
		logger.Fatal("获取响应文本失败: %s", e.Error())
	}

	logger.Info("Hijack结果: %s", text)
}

func CustomApi(id int64) *gpixiv.PixivApi {
	a := gpixiv.NewAPI("GET", fmt.Sprintf("/ajax/illust/%d/pages", id)).
		SetHeader("Accept", "application/json; charset=utf-8").
		SetValue("lang", "en").
		// RespHijack 拦截响应
		SetRespHijack(func(resp *gpixiv.PixivResponse, setBody func(body []byte)) error {
			logger.Info("原始响应状态: %s", resp.Status)

			setBody([]byte("Hijacked"))
			return nil
		})

	return a
}
