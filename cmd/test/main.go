package main

import (
	"github.com/luopengift/golibs/logger"
	"github.com/luopengift/transport"
	_ "github.com/luopengift/transport/api"
	"github.com/luopengift/transport/config"
	"github.com/luopengift/transport/filter"
	_ "github.com/luopengift/transport/plugins"
	"runtime"
	//"flag"
)

const (
	VERSION = "0.0.1"
)

var t *transport.Transport

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	logger.Info("Transport starting...")
	cfg := config.NewConfig()
	logger.Info("%#v", cfg)
	logger.Info("%#v,%#v", transport.InputPlugins, transport.OutputPlugins)


	input := cfg.Input()
	output := cfg.Output()
	logger.Debug("%#v,%#v",input,output)
	//t = transport.NewTransport(input, filter.AddEnter, output)
	t = transport.NewTransport(input, &filter.DefaultConnection{}, output)
	t.Run()
	logger.Debug("%#v", t)

}
