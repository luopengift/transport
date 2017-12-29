package codec

import (
	"github.com/luopengift/transport"
	"github.com/luopengift/types"
)

type Alert struct {
	Labels	map[string]string		`json:"labels"`
	Annotations	map[string]string	`json:"annotations"`
	StartsAt	string			`json:"startsAt,omitempty"`
	EndsAt		string			`json:"endsAt,omitempty"`
	GeneratorURL	string			`json:"generatorURL,omitempty"`
}


// add a enter symbol at end of line, classic written into file
type PrometheusAlertHandler struct{}

func (h *PrometheusAlertHandler) Init(config transport.Configer) error {
	return nil
}

func (h *PrometheusAlertHandler) Handle(in, out []byte) (int, error) {
	src, err := types.ToMap(in)
	if err != nil {
		return 0, err
	}
	labels := map[string]string{
		"alertname": "ERROR_LOG",
		"service": src["serviceName"].(string),//"error_log"
		"file":src["source"].(string),
		"host": src["beat"].(map[string]interface{})["hostname"].(string),
	}
	annotations := map[string]string{
		"summary": src["message"].(string),
	}
	for k, v := range src["error_stack"].(map[string]interface{}) {
		annotations[k] = v.(string)
	}
	dest := Alert{
		Labels: labels,
		Annotations: annotations,
		StartsAt: src["@timestamp"].(string),
		//EndsAt: "0001-01-01T00:00:00Z",
		GeneratorURL: "",
	}
	alerts := []*Alert{&dest}
	b, err := types.ToBytes(alerts)
	if err != nil {
		return 0, nil
	}
	n := copy(out, b)
	return n, nil
}

func (h *PrometheusAlertHandler) Version() string {
	return "0.0.1"
}

func init() {
	transport.RegistHandler("prometheusalert", new(PrometheusAlertHandler))

}
