package transport

import (
	"fmt"
	"io"
)

type Transport struct {
	Inputer
	Outputer
	Buffer []byte
}

func NewTransport(in Inputer, out Outputer) *Transport {
	transport := new(Transport)
	transport.Inputer = in
	transport.Outputer = out
	transport.Buffer = make([]byte, 1024)
	return transport
}

func (t *Transport) Run() {
	for {
		_, err := t.Inputer.Read(t.Buffer)
		switch {
		case err == nil:
			t.Outputer.Write(t.Buffer)
		case err != nil && err == io.EOF:
			fmt.Println("end of inputer")
		default:
			fmt.Println("error:", err)

		}

	}
}
