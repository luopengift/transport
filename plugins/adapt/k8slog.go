package codec

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/luopengift/transport"
	"github.com/luopengift/types"
)

// K8sLogHandler add a enter symbol at end of line, classic written into file
type K8sLogHandler struct {
	GeneratorURL string `json:"generatorURL"`
	CmdbLink     string `json:"cmdbLink"`
	ServiceLink  string `json:"serviceLink"`

	cmdb    []interface{}
	service []interface{}
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
	// 				0	1	2		3							4			5
	// source -> /data/log/mm-prod/edb-channel-9478fc5db-8h8j4/edb-channel/common-error.log
	sources := strings.Split(log.Source, "/")

	log.Namespace = sources[2]
	log.Pod = sources[3]
	log.App = sources[4]

	b, err := types.ToBytes(log)
	if err != nil {
		return 0, fmt.Errorf("%v => %v", err, string(in))
	}
	n := copy(out, b)
	return n, nil
}

// Version version
func (h *K8sLogHandler) Version() string {
	return "0.0.1_091518"
}

func init() {
	transport.RegistHandler("k8slog", new(K8sLogHandler))

}
