package elasticsearch

import (
	"bytes"
	"github.com/luopengift/gohttp"
	"github.com/luopengift/golibs/file"
	"github.com/luopengift/golibs/logger"
	"github.com/luopengift/transport"
)

type EsOutput struct {
	Addrs   []string `json:"addrs"` //es addrs
	Index   string   `json:"index"` //es index
	Type    string   `json:"type"`  //es type
	Timeout int      `json:"time"`  //Pool timeout
	Batch   int      `json:"batch"` //多少条数据提交一次

	Pool   *gohttp.ClientPool
	Buffer chan []byte
}

func NewEsOutput() *EsOutput {
	return new(EsOutput)
}

func (out *EsOutput) Init(config transport.Configer) error {
	out.Timeout = 5
	out.Batch = 1
	err := config.Parse(out)
	if err != nil {
		return err
	}
	out.Pool = gohttp.NewClientPool(5, 50, out.Timeout)
	out.Buffer = make(chan []byte, out.Batch)
	return err
}

func (out *EsOutput) Write(p []byte) (int, error) {
	bulk := NewBulkIndex(file.TimeRule.Handle(out.Index), out.Type, "", p)
	b, err := bulk.Bytes()
	if err != nil {
		return 0, err
	}
	out.Buffer <- b
	return len(b), nil
}

func (out *EsOutput) Start() error {
	cnt := 0
	var buf bytes.Buffer
	for b := range out.Buffer {
		if cnt == 2 {
			logger.Info("send:%v", string(buf.Bytes()))
			err := Send(out.Addrs[0], buf.Bytes())
			if err != nil {
				logger.Error("send bulk error,%v", err)
			}
			buf.Reset()
			cnt = 0
		}
		buf.Write(b)
		cnt = cnt + 1
	}
	return nil
}

func (out *EsOutput) Close() error {
	return nil
}

func init() {
	transport.RegistOutputer("elasticsearch", NewEsOutput())
}
