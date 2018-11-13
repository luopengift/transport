package codec

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/luopengift/log"
	"github.com/luopengift/transport"
)

type Reformat struct {
	Regex string `json:"regex"`
}

type Filebeat struct {
	Timestamp  string            `json:"@timestamp"`
	MetaData   map[string]string `json:"@metadata"`
	Fields     map[string]string `json:"fields"`
	Beat       map[string]string `json:"beat"`
	Host       map[string]string `json:"host"`
	Message    string            `json:"message"`
	Source     string            `json:"source"`
	Offset     int               `json:"offset"`
	Prospector map[string]string `json:"prospector"`
	Input      map[string]string `json:"input"`
}

func (r *Reformat) Init(config transport.Configer) error {
	return config.Parse(r)
}

func (r *Reformat) Handle(in, out []byte) (int, error) {
	log.Info("%v", string(in))
	msg := Format(in)
	b, err := json.Marshal(msg)
	if err != nil {
		return 0, err
	}
	n := copy(out, b)
	log.Info("send %v", string(out))
	return n, nil

}

func (r *Reformat) Version() string {
	return "0.0.1_debug"
}

func init() {
	transport.RegistHandler("reformat", new(Reformat))
}

func Format(msg []byte) map[string]string {
	var err error
	result := map[string]string{}
	recv := &Filebeat{}
	if err = json.Unmarshal(msg, recv); err != nil {
		log.Error("%v", err)
		return nil
	}
	result["topic"] = recv.MetaData["topic"]
	result["src"] = recv.Message

	if result["env"], result["app"], err = formatSource(recv.Source); err != nil {
		log.Warn("%v", err)
	}
	for k, v := range formatMsg(recv.Message) {
		result[k] = v
	}
	return result
}

// /data/log/prod/underwriter-8c87cdc7d-hrfpj/underwriter/underwriter-default.log
func formatSource(source string) (string, string, error) {
	sourceList := strings.Split(source[1:], "/")
	if len(sourceList) != 6 {
		return "", "", log.Errorf("parse source error: %v", source)
	}
	return sourceList[2], sourceList[4], nil
}

var patten = `^(?P<Time>20[0-9]{2}-[0-1][0-9]-[0-3][0-9]\s[0-9]{2}:[0-9]{2}:[0-9]{2},[0-9]{3}Z?)\s\[(?P<traceid>.*?)\s(?P<requestid>.*?)\]\s\[(?P<session>.*?)\]\s+?\[(?P<thread>.*?)\]\s+?(?P<Level>.*?)\s(?P<class>.*?)\s+?-\s(?P<others>.*)$`

func formatMsg(msg string) map[string]string {
	regx := regexp.MustCompile(patten)

	if matchs := regx.FindStringSubmatch(msg); matchs != nil {
		match := make(map[string]string, len(matchs))
		for idx, value := range regx.SubexpNames() {
			match[value] = matchs[idx]
		}
		delete(match, "")
		return match
	}
	return nil
}
