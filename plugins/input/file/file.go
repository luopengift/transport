package file

import (
	"github.com/luopengift/golibs/file"
	"github.com/luopengift/log"
	"github.com/luopengift/transport"
)

const (
	VERSION = "0.0.1"
)

type FileInput struct {
	Path    string `json:"path"`
	EndStop bool   `json:"endstop"`
	*file.Tail
}

func NewFileInput() *FileInput {
	return new(FileInput)
}

func (in *FileInput) Init(config transport.Configer) error {
	err := config.Parse(in)
	if err != nil {
		log.Error("parse error:%v", err)
		return err
	}

	in.Tail = file.NewTail(in.Path, file.TimeRule)
	if in.EndStop {
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

func (in *FileInput) Version() string {
	return VERSION
}

func init() {
	transport.RegistInputer("file", NewFileInput())
}
