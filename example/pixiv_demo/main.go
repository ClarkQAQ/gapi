package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"time"
	"utilware/logger"

	"github.com/ClarkQAQ/gapi"
	"github.com/ClarkQAQ/gapi/apis/pixiv"

	"utilware/gjson"
)

func main() {
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

	// 鉴权 (没有登录无法查看R18作品)

	// (已废弃) 现能获取postKey, 但是无法获取CaptchaID, 后面有时间在看看吧...
	// sessid ,e := p.Login("xxx", "xxx")
	// if e != nil {
	// 	logger.Fatal("登录失败: %s", e.Error())
	// }

	// 手动设置网页SESSID
	// 登录进pixiv.net然后F12到应用或者存储页面获取SESSID
	// 这里是读取环境变量的方式, 你可以直接调用p.SetPHPSESSID("xxx")来设置SESSID
	resp, e := p.Do(pixiv.CookieLogin(os.Getenv("PIXIV_PHPSESSID")))
	if e != nil {
		logger.Fatal("登录失败: %s", e.Error())
	}

	logger.Info("登录成功: %v", string(resp.Raw()))

	// 获取账户的关注动态, 并过滤仅R18的动态
	resp, e = p.Do(pixiv.FollowIllust(1, "r18", ""))
	if e != nil {
		logger.Fatal("获取账户的关注动态失败: %s", e.Error())
	}

	// 获取返回的json数据
	// 并解析写入gjson
	// 除了resp.Raw()其他的返回值都或自动处理gzip压缩
	res, e := resp.GJSON()
	if e != nil {
		logger.Fatal("解析json失败: %s", e.Error())
	}

	fmt.Print("\n\n")
	// 如果觉得慢可以套一层协程, 这里只是为了方便观察就用了单线程处理
	// 多线程记得控制并发数量, 不然会出现超时或者被封IP的情况
	res.Get("body.thumbnails.illust").ForEach(func(key, value gjson.Result) bool {
		artworkId := value.Get("id").Int()
		artworkTitle := value.Get("title").String()
		artworkTags := value.Get("tags").Array()
		artworkAuthor := value.Get("userName").String()
		artworkPageCount := value.Get("pageCount").Int()

		logger.Info("作品编号: %d, 作品标题: %s", artworkId, artworkTitle)
		logger.Info("作品标签: %v", artworkTags)
		logger.Info("作者: %s", artworkAuthor)
		logger.Info("图片数: %d", artworkPageCount)

		// 获取作品详细的图片列表
		resp, e := p.Do(pixiv.GetIllust(artworkId, ""))
		if e != nil {
			logger.Warn("获取作品编号: %d 详细的图片列表失败: %s", artworkId, e.Error())
			return true
		}

		// 同样解析json 获取图片列表
		res, e := resp.GJSON()
		if e != nil {
			fmt.Println(resp.Text())
			fmt.Println("error", e)
			return false
		}

		res.Get("body").ForEach(func(key, value gjson.Result) bool {
			// 原图图片地址
			artworkPicUrl := value.Get("urls.original").String()

			u, e := url.Parse(artworkPicUrl)
			if e != nil {
				logger.Fatal("编号: %d 原始图片URL: %s 解析URL失败(请留意接口数据变化): %s", artworkId, artworkPicUrl, e.Error())
				return true
			}

			logger.Info("编号: %d 原图图片地址: %s", artworkId, u.String())

			// 调用下载图片的函数下载图片
			b, e := p.Do(pixiv.Pximg(u.String()))
			if e != nil {
				logger.Warn("编号: %d 图片URL: %s 下载失败: %s", artworkId, u.String(), e.Error())
				return true
			}

			// 为每个artwork单独保存一个文件夹, 然后生成文件名
			picPath := filepath.Join("test", fmt.Sprintf("%d_%s", artworkId, filepath.Base(u.Path)))
			// 创建文件夹, 在Windows可能会触发UAC
			if e := os.MkdirAll(filepath.Dir(picPath), os.ModePerm); e != nil {
				logger.Fatal("创建文件夹失败: %s", e.Error())
			}

			if e := ioutil.WriteFile(picPath, b.Raw(), os.ModePerm); e != nil {
				logger.Fatal("编号: %d 图片URL: %s 保存失败(请确认是否有权限或者其他问题): %s", artworkId, artworkPicUrl, e.Error())
				return true
			}

			logger.Info("编号: %d 图片URL: %s 下载成功", artworkId, artworkPicUrl)
			return true
		})

		fmt.Print("\n\n")
		return true
	})
}
