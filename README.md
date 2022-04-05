<p align="center">
  <a href="https://github.com/ClarkQAQ/gpixiv">
    <img src="images/logo.png" alt="Logo" width="150" height="80">
  </a>

  <h3 align="center">Gpixiv</h3>
  <p align="center">
    一个动态插件化可任意扩展的 <code>pixiv.net</code> Golang 通用API框架
    <br />
  </p>
</p>

## 目录

- [上手指南](#上手指南)
- [使用到的框架](#使用到的框架)

### 上手指南


##### 一百行代码的一步到位演示文件：[Demo](https://github.com/ClarkQAQ/gpixiv/tree/master/example/demo)

简单的使用方法：

```go
p, e := gpixiv.New(&gpixiv.Options{
	// Pixiv主站地址 不传默认为https://www.pixiv.net
	// 这里是方便测试或者某些使用镜像站点的情况
	URL: "https://www.pixiv.net",
	// 国内特供代理设置 例如: socks5://127.0.0.1:7891
	// 如果有帐号密码需要使用BasicAuth, 例如: socks5://admin:admin@127.0.0.1:7891
	ProxyURL: "socks5://127.0.0.1:7891",
	// 用户代理 不传默认为"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36"
	UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36",
	// 语言 不传默认为"zh"
	// 可选值: "zh", "en", 其他值可以去官网查看
	Language: "zh",
	// 超时时间 不传默认为15秒
	Timeout: 15 * time.Second,
})
if e != nil {
	logger.Fatal("创建Pixiv客户端失败: %s", e.Error())
}

p.SetPHPSESSID(os.Getenv("PIXIV_PHPSESSID"))

resp, e := p.Do(api.FollowIllust(1, "r18", p.Language()))
if e != nil {
	logger.Fatal("获取账户的关注动态失败: %s", e.Error())
}

// 然后resp就是返回的内容了
```

### 使用到的框架

- [gjson](https://github.com/tidwall/gjson)
- [utilware/logger](https://github.com/ClarkQAQ/utilware)

### 版权说明

该项目签署了MIT 授权许可，详情请参阅 [LICENSE.txt](https://github.com/shaojintian/Best_README_template/blob/master/LICENSE.txt)




