package elasticsearch

import (
	"github.com/luopengift/transport"
)

type EsOutput struct {
}


func NewEsOutput() *EsOutput {
    es := new(EsOutput)
    return es
}


func (es *EsOutput) Init(map[string]string) error {
	return nil
}

func (es *EsOutput) Start() error {
    return nil 
}

func (es *EsOutput) Write(p []byte) (int,error) {
	return 0,nil
}

func (es *EsOutput) Close() error {
	return nil
}

func init() {
	transport.RegistOutputer("elasticsearch", NewEsOutput())
}
