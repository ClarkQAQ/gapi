### Pixiv API

> 一个画师常驻的创作网站，提供了一个简单的API，可以获取画师的插画和漫画作品。

#### Pixiv一百行代码的一步到位演示文件：[Demo](https://github.com/ClarkQAQ/gapi/tree/master/example/pixiv_demo)

0. 初始化客户端

```golang
p, e := gapi.New(ccoo.URL, &gapi.Options{
	// ProxyURL: "socks5://127.0.0.1:7891",
})

if e != nil {
	logger.Fatal("创建客户端失败: %s", e.Error())
}

p.SetGHeader(pixiv.GlobalHeader)
```
