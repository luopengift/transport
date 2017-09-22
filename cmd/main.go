package main

import (
	"flag"
	"fmt"
	"github.com/luopengift/golibs/logger"
	"github.com/luopengift/transport"
	"github.com/luopengift/transport/api"
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

func main() {

	cmd := ParseCmdLine()
	
    if cmd.Version {
		fmt.Println("version is", transport.VERSION)
		os.Exit(0)
	}

	if cmd.List {
		fmt.Println(transport.PluginDetail())
		os.Exit(0)
	}

	if cmd.Config == "" {
		logger.Error("config is null,exit...,please -h see help")
		os.Exit(-1)
	}

	cfg := transport.NewConfig(cmd.Config)
	if cfg.Runtime.VERSION != transport.VERSION {
		logger.Warn("runtime version is %s,but config version is %s,NOT match!exit...", transport.VERSION, cfg.Runtime.VERSION)
	}

	if cmd.Read {
		fmt.Println(cfg)
		os.Exit(0)
	}

	logger.Info("Transport starting...")
	runtime.GOMAXPROCS(runtime.NumCPU())

	if cfg.Runtime.DEBUG {
		DebugProfile()
	}
	var err error
	transport.T, err = transport.NewTransport(cfg)
	if err != nil {
		logger.Error("%v", err)
		return
	}
	defer transport.T.Stop()

	api.ApiHttp(cfg.Runtime.HTTP)
	transport.T.Run()
	select {}
}

type CmdLine struct {
    Version bool
    Config  string
    Read    bool
    List    bool
    Pprof   bool
}


func ParseCmdLine() *CmdLine {
    cmdline := new(CmdLine)
	cmdline.Version = *flag.Bool("v", false, "(version)版本号")
	cmdline.Config = *flag.String("f", "", "(config)配置文件")
	cmdline.Read = *flag.Bool("r", false, "(read)读取当前配置文件")
	cmdline.List = *flag.Bool("l", false, "(list)查看插件列表和插件版本")
	cmdline.Pprof = *flag.Bool("p", false, "(pprof)性能")

	flag.Parse()

    return cmdline
}

func DebugProfile() {
	go func() {
		http.ListenAndServe("0.0.0.0:6060", nil)
	}()

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
