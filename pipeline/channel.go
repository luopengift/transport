package pipeline

import (
	"time"
)

type ByteChannel struct {
	channel chan []byte
	closed  chan bool
}

func NewByteChannel(max int) *ByteChannel {
	return &ByteChannel{
		channel: make(chan []byte, max),
		closed:  make(chan bool),
	}
}

func (bc *ByteChannel) Channel() chan []byte {
	return bc.channel
}

func (bc *ByteChannel) Put(b []byte) {
	bc.channel <- b
}

func (bc *ByteChannel) Get() []byte {
	return <-bc.channel
}

func (bc *ByteChannel) Len() int {
	return len(bc.channel)
}

func (bc *ByteChannel) Close() error {
	for {
		if len(bc.channel) > 0 {
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}
	close(bc.channel)
	close(bc.closed)
	return nil

}
