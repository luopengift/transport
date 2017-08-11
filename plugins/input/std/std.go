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

func (std *StdInput) Start() error {
	return nil
}

func (std *StdInput) Close() error {
	return std.Close()
}

func (std *StdInput) Init(config pipeline.Configer) error {
	std.File = os.Stdin
	return nil
}

func init() {
	pipeline.RegistInputer("std", NewStdInput())
}
