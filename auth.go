package gpixiv

import (
	"context"
	"errors"
	"net/http"
)

// 判断是否登录 (需要PHPSESSID)
// 原理是访问用户设置页面, 如果没有302到登录页面, 则已登录
// 当然还有一个判断是否404的分支
func (p *Pixiv) IsLogged() (ret bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	resp, err := p.Request(ctx, http.MethodHead, p.EndpointPATH("/setting_user.php", nil).String(), func(c *http.Client, req *http.Request) error {
		c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}

		return nil
	})

	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusFound {
		return false, err
	} else if resp.StatusCode == http.StatusOK {
		return true, err
	}

	return false, errors.New("pixiv: client: unexpected response for login test")
}

// 设置 PHPSESSID
// 暂时没有自动登录的功能, 只能手动设置
func (p *Pixiv) SetPHPSESSID(sessionID string) {
	p.clearCookies()

	p.c.Jar.SetCookies(
		p.pixivURL,
		[]*http.Cookie{{
			Domain: "." + p.pixivURL.Host,
			Path:   "/",
			Name:   "PHPSESSID",
			Value:  sessionID,
		}},
	)
}
