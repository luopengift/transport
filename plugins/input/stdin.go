package input

import (
	"github.com/luopengift/transport/pipeline"
	"os"
)

type Stdin struct {
	*os.File
}

func NewStdin() *Stdin {
	return new(Stdin)
}

func (stdin *Stdin) Start() error {
	return nil
}

func (stdin *Stdin) Close() error {
	return stdin.Close()
}

func (stdin *Stdin) Init(config pipeline.Configer) error {
	stdin.File = os.Stdin
	return nil
}

func init() {
	pipeline.RegistInputer("stdin", NewStdin())
}
