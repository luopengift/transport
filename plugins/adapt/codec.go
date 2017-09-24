package codec

import (
	"github.com/luopengift/transport"
)

// add a enter symbol at end of line, classic written into file
type AddEnterHandler struct{}

func (h *AddEnterHandler) Init(config transport.Configer) error {
	return nil
}

func (h *AddEnterHandler) Handle(in, out []byte) (int, error) {
	in = append(in, '\n')
	n := copy(out, in)
	return n, nil
}

func (h *AddEnterHandler) Version() string {
	return "0.0.1"
}

// direct connect input and output, do nothing
type NullHandler struct{}

func (h *NullHandler) Init(config transport.Configer) error {
	return nil
}
func (h *NullHandler) Handle(in, out []byte) (int, error) {
	n := copy(out, in)
	return n, nil
}

func (h *NullHandler) Version() string {
	return "0.0.1"
}

func init() {
	transport.RegistHandler("null", new(NullHandler))
	transport.RegistHandler("addenter", new(AddEnterHandler))

}
