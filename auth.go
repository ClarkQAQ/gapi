package gpixiv

import (
	"errors"
	"net/http"
)

func (p *Pixiv) IsLoggedIn() (ret bool, err error) {
	resp, err := p.c.Head(p.EndpointURL("/setting_user.php", nil).String())
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
