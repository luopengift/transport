package elasticsearch

import (
	"github.com/luopengift/golibs/file"
	"github.com/luopengift/transport"
	"os"
)

type EsOutput struct {
	Addrs []string `json:"addrs"` //es addrs
	Index string   `json:"index"` //es index
	Type  string   `json:"type"`  //es type
}

func NewEsOutput() *EsOutput {
	return new(EsOutput)
}

func (out *EsOutput) Init(config transport.Configer) error {
	err := config.Parse(out)
	if err != nil {
		return err
	}
	return err
}

func (out *EsOutput) Write(p []byte) (int, error) {
	return 0, nil
}

func (out *EsOutput) Start() error {
	return nil
}

func (out *EsOutput) Close() error {
	return nil
}

func init() {
	transport.RegistOutputer("elasticsearch", NewEsOutput())
}
