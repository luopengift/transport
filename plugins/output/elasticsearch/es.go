package elasticsearch

import (
	"bytes"
	"github.com/luopengift/gohttp"
	"github.com/luopengift/golibs/file"
	"github.com/luopengift/golibs/logger"
	"github.com/luopengift/transport"
)

const (
	VERSION = "0.0.1"
)

type EsOutput struct {
	Addrs    []string `json:"addrs"`     //es addrs
	Index    string   `json:"index"`     //es index
	Type     string   `json:"type"`      //es type
	Timeout  int      `json:"time"`      //Pool timeout
	Batch    int      `json:"batch"`     //多少条数据提交一次
	MaxProcs int      `json:"max_procs"` //最大并发写协程

	buffer chan []byte
	pool   *gohttp.ClientPool
}

func NewEsOutput() *EsOutput {
	return new(EsOutput)
}

func (out *EsOutput) Init(config transport.Configer) error {
	out.Timeout = 5
	out.Batch = 10
	out.MaxProcs = 10
	err := config.Parse(out)
	if err != nil {
		return err
	}

	out.buffer = make(chan []byte, out.Batch*2)
	// 连接es
	out.pool = gohttp.NewClientPool(out.MaxProcs, out.MaxProcs, out.Timeout)
	return nil
}

func (out *EsOutput) Write(p []byte) (int, error) {
	out.buffer <- p
	return len(p), nil
}

func (out *EsOutput) Start() error {
	for {
		var buf bytes.Buffer
		for tmp := out.Batch; tmp > 0; tmp-- {
			b, ok := <-out.buffer
			if !ok {
				logger.Error("buffer closed")
				return nil
			}
			bulk := NewBulkIndex(file.TimeRule.Handle(out.Index), out.Type, "", b)
			bt, err := bulk.Bytes()
			if err != nil {
				logger.Error("bulk Bytes error:%v", err)
				continue
			}
			buf.Write(bt)
		}
		client, err := out.pool.Get()
		if err != nil {
			logger.Error("get client from pool error:%v", err)
		}
		response, err := client.Url("http://"+out.Addrs[0]).Path("/_bulk").Header("Accept", "application/json").Body(buf.Bytes()).Post()
		if err != nil {
			logger.Error("%v,%v", err, response)
		}
		out.pool.Put(client)
	}
}

func (out *EsOutput) Close() error {
	close(out.buffer)
	return nil
}

func (out *EsOutput) Version() string {
	return VERSION
}

func init() {
	transport.RegistOutputer("elasticsearch", NewEsOutput())
}
