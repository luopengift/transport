package null

import (
	"github.com/luopengift/transport/pipeline"
)

type NullOutput struct {
}

func NewNullOutput() *NullOutput {
	return new(NullOutput)
}

func (n *NullOutput) Init(config pipeline.Configer) error {
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

func init() {
	pipeline.RegistOutputer("null", NewNullOutput())
}
