package codec

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/luopengift/golibs/email"
	"github.com/luopengift/log"
	"github.com/luopengift/transport"
	"github.com/luopengift/types"
)

// K8sLogHandler add a enter symbol at end of line, classic written into file
type K8sLogHandler struct {
	*email.Email
	Subject string `json:"subject" yaml:"subject"`
	From    string `json:"from"`
	To      string `json:"to"`
}

// Init init
func (h *K8sLogHandler) Init(config transport.Configer) error {
	return config.Parse(h)
}

// Log log
type Log struct {
	Message   string `json:"message"`
	Source    string `json:"source"`
	Namespace string `json:"namespace"`
	Pod       string `json:"pod"`
	App       string `json:"app"`
}

// Handle handler
func (h *K8sLogHandler) Handle(in, out []byte) (int, error) {
	msg := &Log{}
	err := json.Unmarshal(in, msg)
	if err != nil {
		return 0, fmt.Errorf("%v => %v", err, string(in))
	}
	sources := strings.Split(msg.Source, "/")

	msg.Namespace = sources[2]
	msg.Pod = sources[3]
	msg.App = sources[4]

	b, err := types.ToBytes(msg)
	if err != nil {
		return 0, fmt.Errorf("%v => %v", err, string(in))
	}

	content := email.NewContent()
	content.Body = msg.Source + "\n\n" + msg.Message
	content.Subject = "mmsz-nx-prod.k8s.local ERROR LOG"

	if err = h.Email.Send(content); err != nil {
		log.Error("send email error: %v", err)
		return 0, err
	}
	log.Info("send email ok!")
	n := copy(out, b)
	return n, err
}

// Version version
func (h *K8sLogHandler) Version() string {
	return "0.0.1_091518"
}

func init() {
	transport.RegistHandler("k8slog", new(K8sLogHandler))

}
