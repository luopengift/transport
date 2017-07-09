package main

import (
	"github.com/luopengift/golibs/logger"
	"github.com/luopengift/transport"
	"github.com/luopengift/transport/config"
	_ "github.com/luopengift/transport/api"
	"github.com/luopengift/transport/filter"
	"github.com/luopengift/transport/plugins/hdfs"
	"github.com/luopengift/transport/plugins/kafka"
	"runtime"
    //"flag"
)

const (
	VERSION = "0.0.1"
)

var output *hdfs.HDFS
var t *transport.Transport

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	logger.Info("Transport starting...")
    cfg := config.NewConfig()
    logger.Info("%#v",cfg)
    
    input := cfg.
    

    t = transport.NewTransport(input, filter.AddEnter, output)
	t.Run()
}
