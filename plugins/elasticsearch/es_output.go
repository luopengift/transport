package elasticsearch

import (
	"github.com/luopengift/transport"
	"fmt"
)

type Elasticsearch struct {}


func (es *Elasticsearch) Init(map[string]string) error {
	return nil
}
func (es *Elasticsearch) Write(p []byte) (int,error) {
	fmt.Println("")
	return 0,nil
}

func (es *Elasticsearch) Close() error {
	return nil
}

func init() {
	transport.RegistOutputer("elasticsearch", new(Elasticsearch))
}
