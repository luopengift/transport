package codec

import (
	"github.com/luopengift/transport"
	"github.com/luopengift/types"
)

type Alert struct {
	Labels	map[string]interface{}		`json:"labels"`
	Annotations	map[string]string	`json:"annotations"`
	StartsAt	string			`json:"startsAt"`
	EndsAt		string			`json:"endsAt"`
	generatorURL	string			`json:"generatorURL"`
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
	labels := map[string]interface{}{
		"type": src["fields"].(map[string]interface{})["service_type"],//"error_log"
		"file":src["source"].(string),
		"host": src["beat"].(map[string]interface{})["hostname"],
	}
	annotations := map[string]string{
		"summary": src["message"].(string),
	}
	dest := Alert{
		Labels: labels,
		Annotations: annotations,
		StartsAt: src["@timestamp"].(string),
		EndsAt: "0001-01-01T00:00:00Z",
	}
	alerts := []*Alert{&dest}
	b, err := types.ToBytes(alerts)
	if err != nil {
		return 0,nil
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
