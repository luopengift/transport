package codec

import (
	"github.com/luopengift/gohttp"
	"github.com/luopengift/transport"
	"strings"
)

type KVHandler struct {
	Keys  []string `json:"keys"`
	Split string   `json:"split"`
}

func (kv *KVHandler) Init(config transport.Configer) error {
	err := config.Parse(kv)
	if err != nil {
		return err
	}
	return err
}

func (kv *KVHandler) Handle(in, out []byte) (int, error) {
	o := make(map[string]string)
	for index, value := range strings.Split(string(in), kv.Split) {
		o[kv.Keys[index]] = value
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
