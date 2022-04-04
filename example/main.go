package main

import (
	"fmt"
	"gpixiv"
	"io/ioutil"
)

func main() {
	p, e := gpixiv.New(&gpixiv.Options{
		ProxyURL: "socks5://127.0.0.1:7891",
	})
	if e != nil {
		panic(e)
	}

	p.SetPHPSESSID("PHPSESSID")

	fmt.Println(p.IsLogged())

	b, e := p.Pximg("https://i.pximg.net/c/250x250_80_a2/custom-thumb/img/2022/03/29/22/23/35/97265461_p0_custom1200.jpg")
	if e != nil {
		panic(e)
	}

	ioutil.WriteFile("test.jpg", b, 0644)
}
