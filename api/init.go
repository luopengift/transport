package api

import (
	"github.com/luopengift/gohttp"
)

func init() {
	app := gohttp.Init()
	go app.Run(":38888")
}
