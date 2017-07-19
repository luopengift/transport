package handler

import (
	"github.com/luopengift/transport"
)

// add a enter symbol at end of line, classic written into file
type AddEnterHandler struct{}

func (h *AddEnterHandler) Handle(in, out []byte) (int, error) {
	in = append(in, '\n')
	n := copy(out, in)
	return n, nil
}

// direct connect input and output, do nothing
type NullHandler struct{}

func (h *NullHandler) Handle(in, out []byte) (int, error) {
	n := copy(out, in)
	return n, nil
}

func init() {
	transport.RegistHandler("null", new(NullHandler))
	transport.RegistHandler("addenter", new(AddEnterHandler))

}
