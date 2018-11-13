package file

import (
	"github.com/luopengift/golibs/file"
	"github.com/luopengift/log"
	"github.com/luopengift/transport"
)

type FilesInput struct {
	Path    []string `json:"path"`
	EndStop bool     `json:"endstop"`

	Files []*file.Tail
	buf   chan []byte
}

func NewFilesInput() *FilesInput {
	return new(FilesInput)
}

func (in *FilesInput) Init(config transport.Configer) error {
	err := config.Parse(in)
	if err != nil {
		log.Error("parse error:%v", err)
		return err
	}
	for _, path := range in.Path {
		tail := file.NewTail(path, file.TimeRule)
		if in.EndStop {
			tail.EndStop(true)
		}
		in.Files = append(in.Files, tail)
	}
	in.buf = make(chan []byte, 1000)
	return nil
}

func (in *FilesInput) Read(p []byte) (int, error) {
	n := copy(p, <-in.buf)
	return n, nil
}

func (in *FilesInput) Start() error {
	for _, tail := range in.Files {
		go func(t *file.Tail) {
			t.ReadLine()
			for msg := range t.NextLine() {
				in.buf <- msg
			}
		}(tail)
	}

	return nil
}

func (in *FilesInput) Close() error {
	for _, tail := range in.Files {
		tail.Close()
	}
	return nil
}

func (in *FilesInput) Version() string {
	return "0.0.1"
}

func init() {
	transport.RegistInputer("files", NewFilesInput())
}
