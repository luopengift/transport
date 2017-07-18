package config

import (
	"github.com/luopengift/golibs/logger"
	"github.com/luopengift/transport"
	"github.com/luopengift/golibs/file"
)


type Configer interface {
	Parse(string) map[string]interface{}
}

type FilterConfig struct{}
type RuntimeConfig struct {
	DEBUG    bool `json:"DEBUG"`
	MAXPROCS int  `json:"MAXPROCS"`
}

func NewRuntimeConfig() *RuntimeConfig {
	return &RuntimeConfig{
		DEBUG:    true,
		MAXPROCS: 1,
	}
}

type InputConfig map[string]string
//type FilterConfig map[string]string
type OutputConfig map[string]string
type HandleConfig map[string]string

type ApiConfig struct {
}



type Config struct {
	Runtime      *RuntimeConfig `json:"runtime"`
	InputConfig  InputConfig   `json:"input"`
	HandleConfig HandleConfig  `json:"handle"`
	OutputConfig OutputConfig  `json:"output"`
	ApiConfig    *ApiConfig     `json:"api"`
}

func NewConfig() *Config {
	cfg := new(Config)
	err := cfg.Init()
	if err != nil {
        logger.Error("config parse error!%v",err)
        return nil
    }
    logger.Warn("Inputer is %#v",cfg.InputConfig)
	logger.Warn("Outputer is %#v",cfg.OutputConfig)
	logger.Warn("Handle is %#v",cfg.HandleConfig)
	return cfg
}

func (cfg *Config) Init() (error) {
    conf := file.NewConfig("./config.json")
    err := conf.Parse(cfg)
    return err

}

func (cfg *Config) Input() transport.Inputer {
	in := transport.InputPlugins[cfg.InputConfig["type"]]
	in.Init(cfg.InputConfig)
	return in
}

func (cfg *Config) Output() transport.Outputer {
    out := transport.OutputPlugins[cfg.OutputConfig["type"]]
    out.Init(cfg.OutputConfig)
	return out
}

func (cfg *Config) Handle() transport.Handler {
    handle := transport.HandlePlugins[cfg.HandleConfig["type"]]
    //handle.Init(cfg.HandleConfig)
	return handle
}

