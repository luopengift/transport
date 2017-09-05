package random

import (
	"github.com/luopengift/golibs/uuid"
	"github.com/luopengift/transport"
    "time"
)

type RandomInput struct {
    Interval    int `json:"interval"`
}

func NewRandomInput() *RandomInput {
	return new(RandomInput)
}

func (in *RandomInput) Start() error {
	return nil
}

func (in *RandomInput) Read(p []byte) (int, error) {
	id := uuid.Rand()
    time.Sleep(time.Duration(in.Interval) * time.Second)
	n := copy(p, id.Hex())
	return n, nil
}

func (in *RandomInput) Close() error {
	return nil
}

func (in *RandomInput) Init(config transport.Configer) error {
    in.Interval = 0
    err := config.Parse(in)
	return err
}

func init() {
	transport.RegistInputer("random", NewRandomInput())
}
