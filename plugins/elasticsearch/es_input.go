package elasticsearch

import (
	"github.com/luopengift/transport"
)

type EsInput struct {
   *ScrollQuery 
    ch chan map[string]interface{}
}

func NewEsInput() *EsInput {
    es := new(EsInput)
    es.ch = make(chan map[string]interface{},10)
    return es
}


func (es *EsInput) Init(cfg map[string]string) error {
    es.ScrollQuery = NewScroll(cfg["protocol"]+"://"+cfg["addrs"],cfg["_index"],cfg["_type"],cfg["scroll"],cfg["query"])
	return nil
}

func (es *EsInput) Start() error {
    es.Next()
    return nil 
}

func (es *EsInput) Read(p []byte) (int,error) {
    return es.ScrollQuery.Read(p)
}

func (es *EsInput) Close() error {
	return nil
}

func init() {
	transport.RegistInputer("elasticsearch", NewEsInput())
}
