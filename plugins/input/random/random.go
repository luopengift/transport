package random

import (
	"github.com/luopengift/golibs/uuid"
	"github.com/luopengift/transport"
)

type RandomInput struct {
}

func NewRandomInput() *RandomInput {
	return new(RandomInput)
}

func (in *RandomInput) Start() error {
	return nil
}

func (in *RandomInput) Read(p []byte) (int, error) {
	id := uuid.Rand()
	n := copy(p, id.Hex())
	return n, nil
}

func (in *RandomInput) Close() error {
	return nil
}

func (in *RandomInput) Init(config transport.Configer) error {
	return nil
}

func init() {
	transport.RegistInputer("random", NewRandomInput())
}
