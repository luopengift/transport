package elasticsearch

import (
	"github.com/luopengift/transport"
)

type Elasticsearch struct {}


func NewEsOutput() *Elasticsearch {
    es := new(Elasticsearch)
    return es
}


func (es *Elasticsearch) Init(map[string]string) error {
	return nil
}

func (es *Elasticsearch) Start() error {
    return nil 
}

func (es *Elasticsearch) Write(p []byte) (int,error) {
	return 0,nil
}

func (es *Elasticsearch) Close() error {
	return nil
}

func init() {
	transport.RegistOutputer("elasticsearch", NewEsOutput())
}
