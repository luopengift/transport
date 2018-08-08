package codec

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/luopengift/gohttp"
	"github.com/luopengift/transport"
	"github.com/luopengift/types"
)

type Alert struct {
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
	StartsAt     time.Time         `json:"startsAt,omitempty"`
	EndsAt       time.Time         `json:"endsAt,omitempty"`
	GeneratorURL string            `json:"generatorURL"`
}

// add a enter symbol at end of line, classic written into file
type PrometheusAlertHandler struct {
	GeneratorURL string `json:"generatorURL"`
	CmdbLink     string `json:"cmdbLink"`
	ServiceLink  string `json:"serviceLink"`

	cmdb    []interface{}
	service []interface{}
}

func (h *PrometheusAlertHandler) QueryCMDB(k, v string) string {
	for _, value := range h.cmdb {
		ec2 := value.(map[string]interface{})
		if tags := ec2["Tags"]; tags != nil {
			for _, tag := range tags.([]interface{}) {
				item := tag.(map[string]interface{})
				if item["Key"].(string) == k && item["Value"].(string) == v {
					return ec2["PrivateIpAddress"].(string)
				}
			}
		} else {
			fmt.Println("Tags 为空: ", ec2["InstanceId"])
		}
	}
	return fmt.Sprintf("IP查询失败:无法查询到%s=%s的资产", k, v)
}

func (h *PrometheusAlertHandler) QuerySERVICE(k, v string) string {
	for _, value := range h.service {
		service := value.(map[string]interface{})
		s, ok := service[k].(string)
		if !ok {
			fmt.Println("service[k].(string) error!", service[k])
			continue
		}

		if s == v {
			if devusers, ok := service["dev_user"].([]interface{}); ok {
				var devs []string
				for _, dev := range devusers {
					devs = append(devs, dev.(string))
				}

				return strings.Join(devs, ",")
			}
			return "devuser为空"
		}
	}
	fmt.Println("service query null", k, v)
	return fmt.Sprintf("开发者查询失败:无法查询到%s=%s的开发者", k, v)
}

// GetExtraInfo get extra info
func GetExtraInfo(url string) ([]interface{}, error) {
	resp, err := gohttp.NewClient().URLString(url).Get()
	if err != nil {
		return nil, err
	}
	if resp.Code()/100 > 2 {
		errMsg := fmt.Sprintf("%s http response is %d", url, resp.Code())
		return nil, fmt.Errorf(errMsg)
	}
	result, err := types.ToMap(resp.Bytes())
	if err != nil {
		return nil, err
	}
	return result["data"].([]interface{}), nil

}

// GetCMDBInfo get cmdb
func GetCMDBInfo(url string) ([]interface{}, error) {
	resp, err := gohttp.NewClient().URLString(url).Get()
	if err != nil {
		return nil, err
	}
	if resp.Code()/100 > 2 {
		errMsg := fmt.Sprintf("%s http response is %d", url, resp.Code())
		return nil, fmt.Errorf(errMsg)
	}

	result, err := types.ToMap(resp.Bytes())
	if err != nil {
		return nil, err
	}
	return result["data"].([]interface{}), nil

}

// Init init
func (h *PrometheusAlertHandler) Init(config transport.Configer) error {
	err := config.Parse(h)
	if err != nil {
		return err
	}
	if h.cmdb, err = GetCMDBInfo(h.CmdbLink); err != nil {
		return err
	}
	fmt.Println(h.cmdb)
	if len(h.cmdb) == 0 {
		return fmt.Errorf("ec2 data is null")
	}
	if h.service, err = GetExtraInfo(h.ServiceLink); err != nil {
		return err
	}
	if len(h.service) == 0 {
		return fmt.Errorf("service data is null")
	}
	fmt.Println(h.service)
	return nil
}

func timeformat(t string, delta int) time.Time {
	zone, _ := time.LoadLocation("Asia/Chongqing")
	loctime, _ := time.ParseInLocation("2006-01-02T15:04:05.000Z", t, zone)
	return loctime.Add(time.Duration(delta) * time.Hour)
}

func getAbstract(msg string) string {
	fristLine, err := bytes.NewBufferString(msg).ReadString('\n')
	if err != nil {
		return ""
	}
	s := strings.Split(fristLine, "ERROR")
	if len(s) == 1 {
		fmt.Println("split by ERROR error!")
		return ""
	}
	ab := strings.Split(s[1], ":")[0]
	if len(ab) > 100 {
		return ab[:100]
	}
	return ab
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

	serviceInfo := h.QuerySERVICE("artifact_id", service)
	fmt.Println("query service result", serviceInfo)

	msg := src["message"].(string)

	labels := map[string]string{
		"alertname": "ERROR_LOG",
		"service":   service,
		"abstract":  getAbstract(msg),
	}
	annotations := map[string]string{
		"summary":    msg,
		"file":       src["source"].(string),
		"host":       host,
		"developers": h.QuerySERVICE("artifact_id", service),
		"ip":         h.QueryCMDB("Name", host), //根据主机名,查询IP地址
	}
	if err_stack, ok := src["error_stack"]; ok {
		for k, v := range err_stack.(map[string]interface{}) {
			annotations[k] = v.(string)
		}
	}

	generatorUrl := h.GeneratorURL
	generatorUrl = strings.Replace(generatorUrl, "SERVICE", service, -1)
	generatorUrl = strings.Replace(generatorUrl, "HOST", host, -1)

	dest := Alert{
		Labels:       labels,
		Annotations:  annotations,
		StartsAt:     timeformat(src["@timestamp"].(string), 8),
		EndsAt:       timeformat(src["@timestamp"].(string), 8+24),
		GeneratorURL: generatorUrl,
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
	return "0.0.9_011118"
}

func init() {
	transport.RegistHandler("prometheusalert", new(PrometheusAlertHandler))

}
