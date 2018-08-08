package hdfs

import (
	"github.com/luopengift/gohttp"
	"github.com/luopengift/golibs/logger"
	"github.com/luopengift/transport"
)

const (
	VERSION = "0.0.2"
)

type HTTPOutput struct {
	Addr string `json:"addr"`
	//Method	string	`json:"method"`
	pool *gohttp.ClientPool
}

func NewHTTPOutput() *HTTPOutput {
	return new(HTTPOutput)
}

func (out *HTTPOutput) Init(config transport.Configer) error {
	err := config.Parse(out)
	if err != nil {
		return err
	}
	out.pool = gohttp.NewClientPool(100, 100, 100)
	return err
}

func (out *HTTPOutput) Start() error {
	return nil
}

func (out *HTTPOutput) Write(p []byte) (int, error) {
	client, err := out.pool.Get()
	if err != nil {
		return 0, err
	}
	defer out.pool.Put(client)
	resp, err := client.KeepAlived(false).URLString(out.Addr).Header("Content-Type", "application/json;charset=utf-8").Body(p).Post()
	//resp, err := gohttp.NewClient().KeepAlived(false).Url(out.Addr).Header("Content-Type", "application/json;charset=utf-8").Body(p).Post()
	if err != nil {
		return 0, err
	}
	if resp.Code() != 200 {
		logger.Error("plugin out http post status not 200, response is: %#v", string(resp.Byte))
		return 0, err
	}
	return len(p), nil
}

func (out *HTTPOutput) Close() error {
	return out.pool.Close()
}
func (out *HTTPOutput) Version() string {
	return VERSION
}

func init() {
	transport.RegistOutputer("http", NewHTTPOutput())
}
