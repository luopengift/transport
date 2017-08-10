package main

import (
	"flag"
	"fmt"
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

/*
Process Signaling

On Unix systems, the NATS server responds to the following signals:

Signal	Result
SIGKILL	Kills the process immediately
SIGINT	Stops the server gracefully
SIGUSR1	Reopens the log file for log rotation
SIGHUP	Reloads server configuration file
*/

const (
	VERSION = "0.0.1"
)

var t *pipeline.Transport

var helpString = `Transport Help
	-h	Help
	-f	config
`

func main() {

	config := flag.String("f", "", "(config)配置文件")
	help := flag.Bool("h", false, "(config)配置文件")
	flag.Parse()

	if *help {
		fmt.Println(helpString)
		return
	}

	if *config == "" {
		fmt.Println(helpString)
		logger.Error("config is null,exit...")
		return
	}

	logger.Info("Transport starting...")
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
