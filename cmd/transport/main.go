package main

import (
	"flag"
	"github.com/luopengift/golibs/logger"
	_ "github.com/luopengift/transport/api"
	"github.com/luopengift/transport/pipeline"
	_ "github.com/luopengift/transport/plugins"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
)

const (
	VERSION = "0.0.1"
)

var t *pipeline.Transport

func main() {

	logger.Info("Transport starting...")
	config := flag.String("f", "", "(config)配置文件")
	flag.Parse()
	if *config == "" {
		logger.Error("config is null,exit...")
		return
	}

	cfg := pipeline.NewConfig(*config)
	logger.Info("%#v", cfg.Runtime)
	runtime.GOMAXPROCS(runtime.NumCPU())

	if cfg.Runtime.DEBUG {
		s := make(chan os.Signal, 1)
		signal.Notify(s, os.Interrupt, os.Kill)

		cpuFile, err := os.OpenFile("./cpu.prof", os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			logger.Error("<file open error> %v", err)
		}
		defer cpuFile.Close()
		pprof.StartCPUProfile(cpuFile)
		defer pprof.StopCPUProfile()
		memFile, err := os.OpenFile("./mem.prof", os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			logger.Error("<file open error> %v", err)
		}
		go func() {
			select {
			case sign := <-s:
				logger.Warn("Get signal:%v, Profile File is cpu.prof/mem.prof", sign)
				//cpu
				pprof.StopCPUProfile()
				cpuFile.Close()
				//mem
				pprof.WriteHeapProfile(memFile)
				memFile.Close()
				os.Exit(-1)
			}
		}()

		go func() {
			http.ListenAndServe("localhost:6060", nil)
		}()
	}

	t = pipeline.NewTransport(cfg)
	defer t.Stop()
	t.Run()

}
