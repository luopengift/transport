package api

import (
    "github.com/luopengift/gohttp"
)

func init() {
    app := gohttp.Init()
    app.Run("38888")
}
