package elasticsearch

import (
	"github.com/luopengift/transport"
)

type EsInput struct {}


func NewEsInput() *EsInput {
    es := new(EsInput)
    return es
}


func (es *EsInput) Init(map[string]string) error {
	return nil
}

func (es *EsInput) Start() error {
    return nil 
}

func (es *EsInput) Read(p []byte) (int,error) {
	return 0,nil
}

func (es *EsInput) Close() error {
	return nil
}

func init() {
	transport.RegistInputer("elasticsearch", NewEsInput())
}
