package elasticsearch

import (
	"github.com/luopengift/golibs/file"
	"github.com/luopengift/transport/pipeline"
)

type EsOutput struct {
	*HttpBulk
	Index string
	Type  string
}

func NewEsOutput() *EsOutput {
	es := new(EsOutput)
	return es
}

func (es *EsOutput) Init(cfg map[string]string) error {
	es.HttpBulk = NewHttpBulk(cfg["protocol"], cfg["addrs"], "", 0, "", "")
	es.Index = cfg["_index"]
	es.Type = cfg["_type"]
	return nil
}

func (es *EsOutput) Start() error {
	return nil
}

func (es *EsOutput) Write(p []byte) (int, error) {
	bulk, err := NewBulkIndex(file.TimeRule.Handle(es.Index), es.Type, "", p).Bytes()
	if err != nil {
		return 0, err
	}
	err = es.HttpBulk.Index(bulk)
	return 0, err
}

func (es *EsOutput) Close() error {
	return nil
}

func init() {
	pipeline.RegistOutputer("elasticsearch", NewEsOutput())
}
