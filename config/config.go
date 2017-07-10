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


type ApiConfig struct {
}



type Config struct {
	Runtime      *RuntimeConfig `json:"runtime"`
	InputConfig  InputConfig   `json:"input"`
	FilterConfig map[string]string  `json:"filter"`
	OutputConfig OutputConfig  `json:"output"`
	ApiConfig    *ApiConfig     `json:"api"`
}

func NewConfig() *Config {
	cfg := new(Config)
	cfg.Init()
	logger.Warn("logger is %#v",cfg.InputConfig)
	logger.Warn("logger is %#v",cfg.OutputConfig)
	return cfg
}

func (cfg *Config) Init() (*Config,error) {
    conf := file.NewConfig("./config.json")
    err := conf.Parse(cfg)
    logger.Info("%+v", conf.String())
    return cfg, err

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



