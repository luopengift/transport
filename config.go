package transport

import (
	"encoding/json"
	"fmt"
	"github.com/luopengift/gohttp"
	"github.com/luopengift/golibs/file"
	"github.com/luopengift/golibs/logger"
	"strings"
)

const (
	B  = 1         //1B = 8bit
	KB = 1024 * B  //1KB
	MB = 1024 * KB //1MB
	GB = 1024 * MB //1GB
	TB = 1024 * GB //1TB
	PB = 1024 * TB //1PB

)

type Configer interface {
	Parse(interface{}) error
}

type RuntimeConfig struct {
	DEBUG    bool   `json:"DEBUG"`
	MAXPROCS int    `json:"MAXPROCS"`
	BYTESIZE int    `json:"BYTESIZE"`
	CHANSIZE int    `json:"CHANSIZE"`
	VERSION  string `json:"VERSION"`
	HTTP     string `json:"HTTP"`
}

func NewRuntimeConfig() *RuntimeConfig {
	return &RuntimeConfig{
		DEBUG:    true,
		MAXPROCS: 1,
		BYTESIZE: 1000,
		CHANSIZE: 100,
		HTTP:     "0.0.0.0:12345",
		VERSION:  "",
	}
}

type pluginConfig map[string]interface{}

func (m pluginConfig) Parse(v interface{}) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, v)
	if err != nil {
		return err
	}
	return nil
}

type Config struct {
	Runtime      *RuntimeConfig          `json:"runtime"`
	InputConfig  map[string]pluginConfig `json:"inputs"`
	HandleConfig map[string]pluginConfig `json:"handles"`
	OutputConfig map[string]pluginConfig `json:"outputs"`
}

func (cfg *Config) String() string {
	var Func = func(cfg map[string]pluginConfig) string {
		str := ""
		writeSpace := " "
		for plugin, config := range cfg {
			str += strings.Repeat(writeSpace, 2) + plugin + ":\n"
			for key, value := range config {
				valueString, _ := gohttp.ToString(value)
				str += strings.Repeat(writeSpace, 4) + key + ": " + valueString + "\n"
			}
		}
		return str
	}
	str := "config info:\n"
	str += "[Inputs]\n"
	str += Func(cfg.InputConfig)
	str += "[Adapts]\n"
	str += Func(cfg.HandleConfig)
	str += "[Outputs]\n"
	str += Func(cfg.OutputConfig)
	return str
}

func NewConfig(path string) *Config {
	cfg := new(Config)
	err := cfg.Init(path)
	if err != nil {
		logger.Error("config parse error!%v", err)
		return nil
	}
	return cfg
}

func (cfg *Config) Init(path string) error {
	conf := file.NewConfig(path)
	err := conf.Parse(cfg)
	return err

}

func (cfg *Config) InitInputs() ([]*Input, error) {
	var inputs []*Input
	for inputName, config := range cfg.InputConfig {
		inputer, ok := Plugins.Inputers[inputName]
		if !ok {
			return nil, fmt.Errorf("[%s] input is not register in pluginspluginConfig", inputName)
		}
		input := NewInput(inputName, inputer)
		if err := input.Inputer.Init(config); err != nil {
			return nil, fmt.Errorf("[%s] input init error:%v", inputName, err)
		}
		inputs = append(inputs, input)
	}
	return inputs, nil
}

func (cfg *Config) InitOutputs() ([]*Output, error) {
	var outputs []*Output
	for outputName, config := range cfg.OutputConfig {
		outputer, ok := Plugins.Outputers[outputName]
		if !ok {
			return nil, fmt.Errorf("[%s] output is not register in pluginspluginConfig", outputName)
		}
		output := NewOutput(outputName, outputer)
		if err := output.Outputer.Init(config); err != nil {
			return nil, fmt.Errorf("[%s] output init error:%v", outputName, err)
		}
		outputs = append(outputs, output)
	}
	return outputs, nil
}

func (cfg *Config) InitCodecs() ([]*Codec, error) {
	var codecs []*Codec
	for codecName, config := range cfg.HandleConfig {
		codec, ok := Plugins.Adapters[codecName]
		if !ok {
			return nil, fmt.Errorf("[%s] codec is not register in pluginspluginConfig", codecName)
		}
		handle := NewCodec(codecName, codec, cfg.Runtime.CHANSIZE)
		if err := handle.Adapter.Init(config); err != nil {
			return nil, fmt.Errorf("[%s] codec init error:%v", codecName, err)
		}
		codecs = append(codecs, handle)
	}
	return codecs, nil
}
