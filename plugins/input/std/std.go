package std

import (
	"github.com/luopengift/transport"
	"os"
)

const (
	VERSION = "0.0.1"
)

type StdInput struct {
	*os.File
}

func NewStdInput() *StdInput {
	return new(StdInput)
}

func (in *StdInput) Start() error {
	return nil
}

func (in *StdInput) Close() error {
	return in.Close()
}

func (in *StdInput) Init(config transport.Configer) error {
	in.File = os.Stdin
	return nil
}

func (in *StdInput) Version() string {
	return VERSION
}

func init() {
	transport.RegistInputer("std", NewStdInput())
}
