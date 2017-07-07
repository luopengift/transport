package config

import (
    "github.com/luopengift/transport"
)

type Configer interface {
   Parse(string) map[string]interface{} 
}

type FilterConfig struct{}
type OutputConfig struct{}
type RuntimeConfig struct {
	DEBUG    bool `json:"DEBUG"`
	MAXPROCS int  `json:"MAXPROCS"`
}
type ApiConfig struct {
}

type Config struct {
	Runtime RuntimeConfig `json:"runtime"`
	Input   InputConfig   `json:"input"`
	Filter  FilterConfig  `json:"filter"`
	Output  OutputConfig  `json:"output"`
	Api     ApiConfig     `json:"api"`
}



type InputConfig struct {
	Type string `json:"type"`
    *transport.Inputer
    Config
}



