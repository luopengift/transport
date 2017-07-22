package main

import (
	"github.com/luopengift/golibs/logger"
	"github.com/luopengift/transport"
	"github.com/luopengift/transport/config"
	_ "github.com/luopengift/transport/plugins"
	"time"
    "runtime"
)

const (
	VERSION = "0.0.1"
)

var t *transport.Transport

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	logger.Info("Transport starting...")
	cfg := config.NewConfig()
	logger.Info("%#v,%#v, %#v", transport.InputPlugins, transport.HandlePlugins, transport.OutputPlugins)

	input := cfg.Input()
	output := cfg.Output()
	handle := cfg.Handle()
	logger.Debug("%#v,%#v,%#v", input, handle, output)
	t = transport.NewTransport(input, handle, output)
	defer t.Stop()
    t.Run()
    select{
        case <- t.End():
            break
    }
    time.Sleep(2 * time.Second)
    return
}
