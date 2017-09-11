package null

import (
	"github.com/luopengift/transport"
)

const (
	VERSION = "0.0.1"
)

type NullOutput struct {
}

func NewNullOutput() *NullOutput {
	return new(NullOutput)
}

func (n *NullOutput) Init(config transport.Configer) error {
	return nil
}

func (n *NullOutput) Write(p []byte) (int, error) {
	return 0, nil
}

func (n *NullOutput) Start() error {
	return nil
}

func (n *NullOutput) Close() error {
	return nil
}

func (n *NullOutput) Version() string {
	return VERSION
}

func init() {
	transport.RegistOutputer("null", NewNullOutput())
}
