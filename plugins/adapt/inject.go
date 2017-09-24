package codec

import (
	"github.com/luopengift/transport"
	"time"
)

// add a enter symbol at end of line, classic written into file
type InjectHandler struct {
	*transport.Inject
}

func (h *InjectHandler) Init(config transport.Configer) error {
	return nil
}

func (h *InjectHandler) Handle(in, out []byte) (int, error) {
	time.Sleep(1 * time.Second) // make program run slow down
	h.InjectInput(in)
	n := copy(out, in)
	return n, nil
}

func (h *InjectHandler) Version() string {
	return "0.0.1_debug"
}

func init() {
	transport.RegistHandler("inject", new(InjectHandler))
}
