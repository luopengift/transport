package codec

import (
	"encoding/json"
	"fmt"
	"github.com/luopengift/gohttp"
	"github.com/luopengift/transport"
	"strconv"
	"strings"
)

type KVHandler struct {
	Keys   [][]string             `json:"keys"`
	Split  string                 `json:"split"`
	Ignore string                 `json:"ignore"`
	Tags   map[string]interface{} `json:"tags"`
}

func (kv *KVHandler) Init(config transport.Configer) error {
	kv.Ignore = "-"
	err := config.Parse(kv)
	if err != nil {
		return err
	}
	return err
}

func (kv *KVHandler) Handle(in, out []byte) (int, error) {
	o := make(map[string]interface{})
	for index, value := range strings.Split(string(in), kv.Split) {
		key := kv.Keys[index][0]
		valueType := kv.Keys[index][1]
		if key == kv.Ignore || valueType == kv.Ignore {
			continue
		}
		switch valueType {
		case "string":
			o[key] = value
		case "int":
			if v, err := strconv.Atoi(value); err == nil {
				o[key] = v
			}
		case "int64":
			if v, err := strconv.ParseInt(value, 10, 64); err == nil {
				o[key] = v
			}
		case "float64":
			if v, err := strconv.ParseFloat(value, 64); err == nil {
				o[key] = v
			}
		case "json":
			v := map[string]interface{}{}
			if err := json.Unmarshal(gohttp.StringToBytes(value), &v); err == nil {
				o[key] = v
			}
		default:
			return 0, fmt.Errorf("type<%v> is unknown", valueType)
		}
	}
	for key, value := range kv.Tags {
		o[key] = value
	}
	b, err := gohttp.ToBytes(o)
	if err != nil {
		return 0, err
	}
	n := copy(out, b)
	return n, nil
}

func init() {
	transport.RegistHandler("kv", new(KVHandler))
}
