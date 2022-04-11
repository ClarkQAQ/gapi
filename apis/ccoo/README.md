### Ccoo

> 写这个只是为了爬鹤山在线数据而已


0. 初始化客户端

```golang
p, e := gapi.New(ccoo.URL, &gapi.Options{
	// ProxyURL: "socks5://127.0.0.1:7891",
})

if e != nil {
	logger.Fatal("创建客户端失败: %s", e.Error())
}

p.SetGHeader(ccoo.GlobalHeader)

// 需要登录的接口需要设置 Cookie
p.Cache().Set("cookie", "COOKIE", 0)

```


1. 获取Job列表

```golang
resp, e := p.Do(ccoo.Jobs("http://www.heshan.ccoo.cn/post/zhaopins/list-0-0-1-0-0-0-1.html"))
if e != nil {
	logger.Fatal("失败: %s", e.Error())
}

data := []ccoo.JobListItem{}
if _, e := resp.JSON(&data); e != nil {
	logger.Fatal("解析失败: %s", e.Error())
}

logger.Info("%+v", data)
```


2. 获取雇主手机号码 (需要登录/Cookie)

```golang

```golang
resp, e := p.Do(ccoo.Phone("http://www.heshan.ccoo.cn/home/recruit/getpositiontel.html", "xxxxxxx"))
if e != nil {
	logger.Fatal("失败: %s", e.Error())
}

logger.Info("%s", resp.Raw())
```

