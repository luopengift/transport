package transport

import (
	"encoding/json"
	"fmt"
	"github.com/luopengift/gohttp"
	"github.com/luopengift/golibs/file"
	"github.com/luopengift/golibs/logger"
	"strings"
)

type Configer interface {
	Parse(interface{}) error
}

type RuntimeConfig struct {
	DEBUG    bool   `json:"DEBUG"`
	MAXPROCS int    `json:"MAXPROCS"`
	VERSION  string `json:"VERSION"`
	HTTP     string `json:"HTTP"`
}

func NewRuntimeConfig() *RuntimeConfig {
	return &RuntimeConfig{
		DEBUG:    true,
		MAXPROCS: 1,
		HTTP:     "0.0.0.0:12345",
		VERSION:  "",
	}
}

type Map map[string]interface{}

func (m Map) Int(key string) int {
	return m[key].(int)
}

func (m Map) Strings(key string) []string {
	var ret []string
	for _, v := range m[key].([]interface{}) {
		ret = append(ret, v.(string))
	}
	return ret
}

func (m Map) String(key string) string {
	return m[key].(string)
}

func (m Map) Parse(v interface{}) error {
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
	Runtime      *RuntimeConfig `json:"runtime"`
	InputConfig  map[string]Map `json:"inputs"`
	HandleConfig map[string]Map `json:"handles"`
	OutputConfig map[string]Map `json:"outputs"`
}

func (cfg *Config) String() string {
	var Func = func(cfg map[string]Map) string {
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
	str += "[Codecs]\n"
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
			return nil, fmt.Errorf("[%s] input is not register in pluginsMap", inputName)
		}
		input := NewInput(inputName, inputer)
		input.Inputer.Init(config)
		inputs = append(inputs, input)
	}
	return inputs, nil
}

func (cfg *Config) InitOutputs() ([]*Output, error) {
	var outputs []*Output
	for outputName, config := range cfg.OutputConfig {
		outputer, ok := Plugins.Outputers[outputName]
		if !ok {
			return nil, fmt.Errorf("[%s] output is not register in pluginsMap", outputName)
		}
		output := NewOutput(outputName, outputer)
		output.Outputer.Init(config)
		outputs = append(outputs, output)
	}
	return outputs, nil
}

func (cfg *Config) InitCodecs() ([]*Codec, error) {
	var handles []*Codec
	for handleName, config := range cfg.HandleConfig {
		handler, ok := Plugins.Handlers[handleName]
		if !ok {
			return nil, fmt.Errorf("[%s] handle is not register in pluginsMap", handleName)
		}
		handle := NewCodec(handleName, handler, 1000)
		handle.Handler.Init(config)
		handles = append(handles, handle)
	}
	return handles, nil
}
