package output

import (
	"github.com/luopengift/transport/pipeline"
	"os"
)

type Stdout struct {
	*os.File
}

func NewStdout() *Stdout {
	return new(Stdout)
}

func (stdout *Stdout) Start() error {
	return nil
}

func (stdout *Stdout) Close() error {
	return stdout.Close()
}

func (stdout *Stdout) Init(config pipeline.Configer) error {
	stdout.File = os.Stdout
	return nil
}

func init() {
	pipeline.RegistOutputer("stdout", NewStdout())
}
