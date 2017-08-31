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
	Timeout int      `json:"time"`

	Pool   *gohttp.ClientPool
	Buffer *bytes.Buffer
}

func NewEsOutput() *EsOutput {
	return new(EsOutput)
}

func (out *EsOutput) Init(config transport.Configer) error {
	out.Timeout = 5
	err := config.Parse(out)
	if err != nil {
		return err
	}
	out.Pool = gohttp.NewClientPool(5, 50, out.Timeout)
	out.Buffer = bytes.NewBuffer(make([]byte, 1*transport.M))
	return err
}

func (out *EsOutput) Write(p []byte) (int, error) {
	bulk := NewBulkIndex(file.TimeRule.Handle(out.Index), out.Type, "", p)
	b, err := bulk.Bytes()
	if err != nil {
		return 0, err
	}
	if out.Buffer.Cap()-out.Buffer.Len() < len(b) {
		err = Send(out.Addrs[0], out.Buffer.Bytes())
		if err != nil {
			logger.Error("send bulk error,%v", err)
		}
		out.Buffer.Reset()
	}
	return out.Buffer.Write(b)
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
