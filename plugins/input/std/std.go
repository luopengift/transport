package std

import (
	"github.com/luopengift/transport/pipeline"
	"os"
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

func (in *StdInput) Init(config pipeline.Configer) error {
	in.File = os.Stdin
	return nil
}

func init() {
	pipeline.RegistInputer("std", NewStdInput())
}
