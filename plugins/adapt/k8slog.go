package codec

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/luopengift/golibs/email"
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
	log := &Log{}
	err := json.Unmarshal(in, log)
	if err != nil {
		return 0, fmt.Errorf("%v => %v", err, string(in))
	}
	sources := strings.Split(log.Source, "/")

	log.Namespace = sources[2]
	log.Pod = sources[3]
	log.App = sources[4]

	b, err := types.ToBytes(log)
	if err != nil {
		return 0, fmt.Errorf("%v => %v", err, string(in))
	}

	content := email.NewContent()
	content.Body = log.Source + "\n\n" + log.Message
	content.Subject = "mmsz-nx-prod.k8s.local ERROR LOG"

	err = h.Email.Send(content)
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
