<p align="center">
  <a href="https://github.com/ClarkQAQ/gapi">
    <img src="images/logo.png" alt="Logo" width="150" height="80">
  </a>

  <h3 align="center">Gapi</h3>
  <p align="center">
    一个动态插件化可任意扩展的 Golang 通用API框架
    <br />
  </p>
</p>


## 目录

- [上手指南](#上手指南)
- [使用到的框架](#使用到的框架)

### 上手指南


#### Pixiv一百行代码的一步到位演示文件：[Demo](https://github.com/ClarkQAQ/gapi/tree/master/example/demo)

##### 简单的使用方法：

```go
p, e := gapi.New(pixiv.URL, &gapi.Options{
	// 国内特供代理设置 例如: socks5://127.0.0.1:7891
	// 如果有帐号密码需要使用BasicAuth, 例如: socks5://admin:admin@127.0.0.1:7891
	ProxyURL: "socks5://127.0.0.1:7891",
	// 超时时间 不传默认为15秒
	Timeout: 15 * time.Second,
})
if e != nil {
	logger.Fatal("创建Gapi客户端失败: %s", e.Error())
}

// 设置api定制全局header
p.SetGHeader(pixiv.GlobalHeader)

resp, e := p.Do(pixiv.CookieLogin(os.Getenv("PIXIV_PHPSESSID")))
if e != nil {
	logger.Fatal("登录失败: %s", e.Error())
}
logger.Info("登录成功: %v", string(resp.Raw()))
```

##### API 自定义插件:

完整测试文件：[Api](https://github.com/ClarkQAQ/gapi/tree/master/example/api)

```go
api.New("GET", fmt.Sprintf("/ajax/illust/%d/pages", 1)).
SetHeader("Accept", "application/json; charset=utf-8").
SetValue("lang", "en").
// RespHijack 拦截响应
SetRespHijack(func(resp *http.Response, respBody func(b []byte) []byte) error {
	logger.Info("原始响应状态: %s", resp.Status)
	respBody([]byte("Hijacked"))
	return nil
})
```

### 使用到的框架

- [gjson](https://github.com/tidwall/gjson)
- [utilware/logger](https://github.com/ClarkQAQ/utilware)

### 版权说明

该项目签署了MIT 授权许可，详情请参阅 [LICENSE.txt](https://github.com/shaojintian/Best_README_template/blob/master/LICENSE.txt)




