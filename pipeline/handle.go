package pipeline

import (
	"github.com/luopengift/golibs/channel"
)

type Handler interface {
	Handle(in, out []byte) (n int, err error)
}

type Middleware struct {
	Handler
	*channel.Channel
}

func NewMiddleware(h Handler, max int64) *Middleware {
	m := new(Middleware)
	m.Handler = h
	m.Channel = channel.NewChannel(max)
	return m
}

func (m *Middleware) Handle(in, out []byte) (int, error) {
	return m.Handler.Handle(in, out)
}
