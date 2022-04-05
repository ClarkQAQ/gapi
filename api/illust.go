package api

import (
	"errors"
	"fmt"
)

// 获取账户的关注动态
// page => 页数
// mode => 模式 (all, r18)
// lang => 语言 (zh, en) 或者调用p.Language()
func FollowIllust(page int, mode string, lang string) *PixivApi {
	return New("GET", "/ajax/follow_latest/illust").
		SetHeader("Accept", "application/json; charset=utf-8").
		SetValue("p", fmt.Sprint(page)).
		SetValue("mode", mode).
		SetValue("lang", lang)
}

// 获取图片列表
// id => 插图ID/漫画ID
// lang => 语言 (zh, en) 或者调用p.Language()
func GetIllust(id int64, lang string) *PixivApi {
	a := New("GET", fmt.Sprintf("/ajax/illust/%d/pages", id)).
		SetHeader("Accept", "application/json; charset=utf-8").
		SetValue("lang", lang)

	if id <= 0 {
		a.SetError(errors.New("artwork id must be greater than 0"))
	}

	return a
}
