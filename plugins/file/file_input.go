package file

import (
	"github.com/luopengift/transport"
	"github.com/luopengift/golibs/file"
)

type Input struct {
	*file.Tail
}

func NewFileInput() *Input {
	return new(Input)
}

func (in *Input) Init(cfg map[string]string) error {
    in.Tail = file.NewTail(cfg["path"],file.NullRule)
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
        return nil//stdout.Close()
}


func init() {
    transport.RegistInputer("file",NewFileInput())
}

