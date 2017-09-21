package transport

import (
	"github.com/luopengift/golibs/channel"
	"sync"
	"sync/atomic"
)

type Adapter interface {
	Init(config Configer) error
	Handle(in, out []byte) (n int, err error)
	Version() string
}

type Adapt struct {
	Name string
	cnt  uint64
	*sync.Mutex
	channel *channel.Channel
	Adapter
}

func NewAdapt(name string, a Adapter, size int) *Adapt {
	at := new(Adapt)
	at.Name = name
	at.cnt = 0
	at.Adapter = a
	at.Mutex = new(sync.Mutex)
	at.channel = channel.NewChannel(size)
	return at
}

func (at *Adapt) Init(config Configer) error {
	return nil
}

func (at *Adapt) Count() uint64 {
	return at.cnt
}
func (at *Adapt) Handle(in, out []byte) (int, error) {
	n, err := at.Adapter.Handle(in, out)
	atomic.AddUint64(&at.cnt, 1)
	return n, err
}
