package std

import (
	"github.com/luopengift/transport/pipeline"
	"os"
)

type StdOutput struct {
	*os.File
}

func NewStdOutput() *StdOutput {
	return new(StdOutput)
}

func (out *StdOutput) Start() error {
	return nil
}

func (out *StdOutput) Close() error {
	return out.Close()
}

func (out *StdOutput) Init(config pipeline.Configer) error {
	out.File = os.Stdout
	return nil
}

func init() {
	pipeline.RegistOutputer("std", NewStdOutput())
}
