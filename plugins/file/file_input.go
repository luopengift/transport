package file

import (
	"github.com/luopengift/golibs/file"
	"github.com/luopengift/transport"
)

type Input struct {
	*file.Tail
}

func NewFileInput() *Input {
	return new(Input)
}

func (in *Input) Init(cfg map[string]string) error {
	in.Tail = file.NewTail(cfg["path"], file.TimeRule)
	return nil
}

func (in *Input) Read(p []byte) (int, error) {
	return in.Tail.Read(p)
}

func (in *Input) Start() error {
	in.Tail.ReadLine()
	return nil
}

func (in *Input) Close() error {
	return in.Tail.Close()
}

func init() {
	transport.RegistInputer("file", NewFileInput())
}
