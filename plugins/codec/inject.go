package codec

import (
	"github.com/luopengift/transport"
	"time"
)

// add a enter symbol at end of line, classic written into file
type DebugInjectHandler struct {
	*transport.Inject
}

func (h *DebugInjectHandler) Init(config transport.Configer) error {
	return nil
}

func (h *DebugInjectHandler) Handle(in, out []byte) (int, error) {
	time.Sleep(1 * time.Second) // make program run slow down
	h.InjectInput(in)
	n := copy(out, in)
	return n, nil
}

func init() {
	transport.RegistHandler("DEBUG_InjectInput", new(DebugInjectHandler))
}
