package elasticsearch

import (
	"github.com/luopengift/transport/pipeline"
)

type EsInput struct {
	*ScrollQuery
}

func NewEsInput() *EsInput {
	es := new(EsInput)
	return es
}

func (es *EsInput) Init(cfg map[string]string) error {
	es.ScrollQuery = NewScroll(cfg["protocol"]+"://"+cfg["addrs"], cfg["_index"], cfg["_type"], cfg["scroll"], cfg["query"])
	return nil
}

func (es *EsInput) Start() error {
	return es.Next()
}

func (es *EsInput) Read(p []byte) (int, error) {
	return es.ScrollQuery.Read(p)
}

func (es *EsInput) Close() error {
	return es.ScrollQuery.Close()
}

func init() {
	pipeline.RegistInputer("elasticsearch", NewEsInput())
}
