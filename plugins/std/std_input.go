package plugins

import (
	"github.com/luopengift/transport"
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

func (stdin *Stdin) Init(map[string]string) error {
	stdin.File = os.Stdin
	return nil
}

func init() {
	transport.RegistInputer("stdin", NewStdin())
}
