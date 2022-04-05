package api

import (
	"fmt"
	"net/http"
	"net/url"
)

// 获取账户的关注动态
// page => 页数
// mode => 模式 (all, r18)
// lang => 语言 (zh, en) 或者调用p.Language()
func FollowIllust(page int, mode string, lang string) *PixivApi {
	return New("GET", "/ajax/follow_latest/illust", http.Header{
		"Accept": []string{"application/json; charset=utf-8"},
	}, url.Values{
		"p":    {fmt.Sprint(page)},
		"mode": {mode},
		"lang": {lang},
	}, nil, nil)
}

// 获取图片列表
// id => 插图ID/漫画ID
// lang => 语言 (zh, en) 或者调用p.Language()
func GetIllust(id int64, lang string) *PixivApi {
	return New("GET", fmt.Sprintf("/ajax/illust/%d/pages", id), http.Header{
		"Accept": []string{"application/json; charset=utf-8"},
	}, url.Values{
		"lang": {lang},
	}, nil, nil)
}
