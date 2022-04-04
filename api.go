package gpixiv

import (
	"gpixiv/api"
)

// Pixiv接口返回
type PixivResponse struct {
	api        *api.PixivApi
	statusCode int
	accept     string
	body       []byte
	error      error
}

// 请求Pixiv接口 (接口需通过api包里面的方法生成)
func (p *Pixiv) Do(api *api.PixivApi) (*PixivResponse, error) {

	return nil, nil
}
