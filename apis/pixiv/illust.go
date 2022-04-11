package pixiv

import (
	"errors"
	"fmt"

	"github.com/ClarkQAQ/gapi"
)

// 获取账户的关注动态
// page => 页数
// mode => 模式 (all, r18)
// lang => 语言 (zh, en)
func FollowIllust(page int, mode string, lang string) *gapi.GapiApi {
	return gapi.NewAPI("GET", "/ajax/follow_latest/illust").
		SetHeader("Accept", "application/json; charset=utf-8").
		SetValue("p", fmt.Sprint(page)).
		SetValue("mode", mode).
		SetValue("lang", lang)
}

// 获取图片列表
// id => 插图ID/漫画ID
// lang => 语言 (zh, en)
func GetIllust(id int64, lang string) *gapi.GapiApi {
	a := gapi.NewAPI("GET", fmt.Sprintf("/ajax/illust/%d/pages", id)).
		SetHeader("Accept", "application/json; charset=utf-8").
		SetValue("lang", lang)

	if id <= 0 {
		a.SetError(errors.New("artwork id must be greater than 0"))
	}

	return a
}
