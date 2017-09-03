package codec

import (
	"fmt"
	"github.com/luopengift/gohttp"
	"github.com/luopengift/transport"
	"regexp"
)

type GrokHandler struct {
	Regex string `json:"regex"`

	Regexp *regexp.Regexp
}

func (j *GrokHandler) Init(config transport.Configer) error {
	err := config.Parse(j)
	if err != nil {
		return err
	}
	j.Regexp, err = regexp.Compile(j.Regex)
	return err
}

func (j *GrokHandler) Handle(in, out []byte) (int, error) {
	match := j.Regexp.FindStringSubmatch(string(in))
	if match == nil {
		return 0, fmt.Errorf("can't regex input %v", string(in))
	}
	kv := make(map[string]string)
	for key, value := range j.Regexp.SubexpNames() {
		kv[value] = match[key]
	}
	delete(kv, "")
	b, err := gohttp.ToBytes(kv)
	if err != nil {
		return 0, err
	}
	n := copy(out, b)
	return n, nil
}

func init() {
	transport.RegistHandler("grok", new(GrokHandler))
}
