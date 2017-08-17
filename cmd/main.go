package main

import (
	"flag"
	"github.com/luopengift/golibs/logger"
	_ "github.com/luopengift/transport/api"
	"github.com/luopengift/transport"
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
	VERSION = "0.0.2"
)

var t *transport.Transport

func main() {

	config := flag.String("f", "", "(config)配置文件")
	flag.Parse()

	if *config == "" {
		logger.Error("config is null,exit...,please -h see help")
		return
	}

	logger.Info("Transport starting...")

	cfg := transport.NewConfig(*config)
	if cfg.Runtime.VERSION != VERSION {
		logger.Error("runtime version is %s,but config version is %s,NOT match!exit...", VERSION, cfg.Runtime.VERSION)
		return
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	if cfg.Runtime.DEBUG {

		DebugProfile()

		go func() {
			http.ListenAndServe("localhost:6060", nil)
		}()
	}
    var err error
	t, err = transport.NewTransport(cfg)
    if err != nil {
        logger.Error("%v",err)
        return
    }
	defer t.Stop()
	t.Run()
	select {}
}

func DebugProfile() {
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt, os.Kill)
	//cpu
	cpu, err := os.Create("./var/cpu.prof")
	if err != nil {
		logger.Error("<file open error> %v", err)
	}
	err = pprof.StartCPUProfile(cpu)
	if err != nil {
		logger.Error("could not start CPU profile:%v", err)
	}
	// memory
	mem, err := os.Create("./var/mem.prof")
	if err != nil {
		logger.Error("<file open error> %v", err)
	}

	go func() {
		select {
		case sign := <-s:
			logger.Warn("Get signal:%v, Profile File is cpu.prof/mem.prof", sign)
			//cpu
			pprof.StopCPUProfile()
			if err := cpu.Close(); err != nil {
				logger.Error("cpu pprof file close err:%v", err)
			}
			//mem
			runtime.GC() // get up-to-date statistics
			if err := pprof.WriteHeapProfile(mem); err != nil {
				logger.Error("could not write memory profile:%v", err)
			}
			if err := mem.Close(); err != nil {
				logger.Error("mem pprof file close err:%v", err)
			}
			os.Exit(-1)
		}
	}()
	logger.Warn("Starting loading performance data, please press CTRL+C exit...")
}
