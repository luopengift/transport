package file

import (
	//	"encoding/json"
	"github.com/luopengift/golibs/file"
	"github.com/luopengift/golibs/logger"
	"github.com/luopengift/transport/pipeline"
)

type FileInput struct {
	*file.Tail
}

func NewFileInput() *FileInput {
	return new(FileInput)
}

type FileConfig struct {
	Path    []string `json:"path"`
	EndStop bool     `json:"endstop"`
}

func (in *FileInput) Init(config pipeline.Configer) error {
	cfg := FileConfig{}
	err := config.Parse(&cfg)
	if err != nil {
		logger.Error("parse error:%v", err)
		return err
	}

	in.Tail = file.NewTail(cfg.Path[0], file.TimeRule)
	if cfg.EndStop {
		in.Tail.EndStop(true)
	}
	return nil
}

func (in *FileInput) Read(p []byte) (int, error) {
	return in.Tail.Read(p)
}

func (in *FileInput) Start() error {
	in.Tail.ReadLine()
	return nil
}

func (in *FileInput) Close() error {
	return in.Tail.Close()
}

func init() {
	pipeline.RegistInputer("file", NewFileInput())
}
