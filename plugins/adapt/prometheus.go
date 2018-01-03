package codec

import (
	"time"
	"fmt"
	"github.com/luopengift/transport"
	"github.com/luopengift/types"
	"strings"
)

type Alert struct {
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
	StartsAt     time.Time         `json:"startsAt,omitempty"`
	EndsAt       time.Time         `json:"endsAt,omitempty"`
	GeneratorURL string            `json:"generatorURL"`
}

// add a enter symbol at end of line, classic written into file
type PrometheusAlertHandler struct{}

func (h *PrometheusAlertHandler) Init(config transport.Configer) error {
	return nil
}

func timeformat(t string) string {
	zone, _ := time.LoadLocation("Asia/Chongqing")
	loctime, _ := time.ParseInLocation("2006-01-02T15:04:05.000Z", t, zone)
	return loctime.Add(8 * time.Hour).Format(time.RFC822)
}

func timeformat2(t string) time.Time {
	zone, _ := time.LoadLocation("Asia/Chongqing")
	loctime, _ := time.ParseInLocation("2006-01-02T15:04:05.000Z", t, zone)
	return loctime.Add(8 * time.Hour)
}

func (h *PrometheusAlertHandler) Handle(in, out []byte) (int, error) {
	src, err := types.ToMap(in)
	if err != nil {
		return 0, fmt.Errorf("%v => %v", err, string(in))
	}
	file := ""
	value, ok := src["source"]
	if ok {
		file = value.(string)
		//return 0, fmt.Errorf("missing source field!")
	}
	host := src["beat"].(map[string]interface{})["hostname"].(string)
	service := ""
	services := strings.Split(file, "/")
	if len(services) < 7 {
		service = file
	} else {
		service = services[7]
	}
	url := "http://10.28.13.66:5601/app/kibana#/discover?_g=(refreshInterval:(display:Off,pause:!f,value:0),time:(from:now-24h,mode:quick,to:now))&_a=(columns:!(_source),index:'logstash-*',interval:auto,query:(query_string:(analyze_wildcard:!t,query:'serviceName:%22"+service+"%22%20AND%20host:%22"+host+"%22%20AND%20level:%22ERROR%22')),sort:!('@timestamp',desc))"
	labels := map[string]string{
		"alertname": "ERROR_LOG",
		"service":   service, //"error_log"
		"file":      src["source"].(string),
		"host":      host,
	//	"starttime": timeformat(src["@timestamp"].(string)),
	}
	annotations := map[string]string{
		"summary": src["message"].(string),
	}
	if err_stack, ok := src["error_stack"]; ok {
		for k, v := range err_stack.(map[string]interface{}) {
			annotations[k] = v.(string)
		}
	}
	dest := Alert{
		Labels:      labels,
		Annotations: annotations,
		StartsAt:    timeformat2(src["@timestamp"].(string)),
		//EndsAt: "0001-01-01T00:00:00Z",
		GeneratorURL: url,
	}
	alerts := []*Alert{&dest}
	b, err := types.ToBytes(alerts)
	if err != nil {
		return 0, fmt.Errorf("%v => %v", err, string(in))
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
