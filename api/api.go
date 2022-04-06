package api

import (
	"errors"
	"fmt"
	"gpixiv"
	"time"
)

func CheckPHPSESSID(p *gpixiv.Pixiv) error {
	for _, v := range p.Client().Jar.Cookies(p.GetURL()) {
		if v.Name == "PHPSESSID" {
			return nil
		}
	}

	return errors.New("this api require PHPSESSID")
}

func GetCsrfTokenString(p *gpixiv.Pixiv) (string, error) {
	if v, ok := p.Cache().Get("csrf_token"); ok {
		return fmt.Sprint(v), nil
	}

	// 获取csrf_token
	resp, e := p.Do(GetCsrfToken())
	if e != nil {
		return "", e
	}

	csrfToken := string(resp.Raw())
	if e != nil {
		return "", errors.New("get csfr_token failed")
	}

	p.Cache().Set("csrf_token", csrfToken, 15*time.Minute)
	return csrfToken, nil
}
