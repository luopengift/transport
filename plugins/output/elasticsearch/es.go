package elasticsearch

import (
	"bytes"
	"github.com/luopengift/gohttp"
	"github.com/luopengift/golibs/file"
	"github.com/luopengift/golibs/logger"
	"github.com/luopengift/transport"
	"os"
)

type EsOutput struct {
	Addrs   []string `json:"addrs"` //es addrs
	Index   string   `json:"index"` //es index
	Type    string   `json:"type"`  //es type
	Timeout int      `json:"time"`

	Pool   *gohttp.ClientPool
	Buffer bytes.Buffer
}

func NewEsOutput() *EsOutput {
	return new(EsOutput)
}

func (out *EsOutput) Init(config transport.Configer) error {
	err := config.Parse(out{
		Timeout: 5,
	})
	if err != nil {
		return err
	}
	out.Pool = gohttp.NewClientPool(5, 50, out.Timeout)
	out.Buffer = bytes.NewBuffer(make([]byte, 1*transport.M))
	return err
}

func (out *EsOutput) Write(p []byte) (int, error) {
	bulk, err := NewBulkIndex(file.TimeRule.Handle(out.Index), out.Type, "", p)
	if err != nil {
		return 0, err
	}
	if out.Buffer.Cap()-out.Buffer.Len() < len(bulk) {
		err = Send(out.Addrs[0], out.Buffer.Bytes())
		if err != nil {
			logger.Error("send bulk error,%v", err)
		}
		out.Buffer.Reset()
	}
	return out.Buffer.Write(p)
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
