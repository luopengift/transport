package pipeline

import (
	"github.com/luopengift/golibs/channel"
)

type Handler interface {
	Init(config Configer) error
	Handle(in, out []byte) (n int, err error)
}

type Codec struct {
	Name string
	Handler
	*channel.Channel
}

func NewCodec(name string, h Handler, max int) *Codec {
	m := new(Codec)
	m.Name = name
	m.Handler = h
	m.Channel = channel.NewChannel(max)
	return m
}

func (m *Codec) Handle(in, out []byte) (int, error) {
	return m.Handler.Handle(in, out)
}
