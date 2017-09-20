package influxdb

import (
	"encoding/json"
	"fmt"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/luopengift/golibs/logger"
	"github.com/luopengift/transport"
	"net/url"
	"time"
)

type data struct {
	Name   string                 `json:"name"`
	Fields map[string]interface{} `json:"fields"`
	Tags   map[string]string      `json:"tags"`
	Time   int64                  `json:"time"`
}

type InfluxOutput struct {
	Addr      string `json:"addr"`
	DB        string `json:"database"`
	Precision string `json:"precision"`
	User      string `json:"username"`
	Pass      string `json:"password"`
	Batch     int    `json:"batch"`

	buffer chan *client.Point
	client client.Client
}

func NewInfluxOutput() *InfluxOutput {
	return new(InfluxOutput)
}

func (out *InfluxOutput) Init(cfg transport.Configer) error {
	out.Batch = 1
	out.DB = "mydb"
	out.Precision = "ns"
	err := cfg.Parse(out)
	if err != nil {
		return err
	}

	u, err := url.Parse(out.Addr)
	if err != nil {
		return err
	}

	switch u.Scheme {
	case "udp":
		out.client, err = client.NewUDPClient(client.UDPConfig{
			Addr: out.Addr,
		})
	case "http", "https":
		out.client, err = client.NewHTTPClient(client.HTTPConfig{
			Addr:     out.Addr,
			Username: out.User,
			Password: out.Pass,
		})
	default:
		return fmt.Errorf("scheme error")
	}
	out.buffer = make(chan *client.Point, out.Batch*2)
	return err
}

func (out *InfluxOutput) Start() error {
	for {
		bp, err := client.NewBatchPoints(client.BatchPointsConfig{
			Database: out.DB,
			//Precision is the write precision of the points, defaults to "ns"
			Precision: out.Precision,
		})
		if err != nil {
			return err
		}

		for tmp := out.Batch; tmp > 0; tmp-- {
			point, ok := <-out.buffer
			if !ok {
				logger.Error("buffer closed")
				return transport.BufferClosedError
			}
			bp.AddPoint(point)
		}
		err = out.client.Write(bp)
		if err != nil {
			logger.Error("write error:%v", err)
		}
		logger.Info("%#v", bp.Points()[0].String())
	}
}

func (out *InfluxOutput) Write(p []byte) (int, error) {
	dat := data{}
	err := json.Unmarshal(p, &dat)
	if err != nil {
		return 0, err
	}
	pt, err := client.NewPoint(dat.Name, dat.Tags, dat.Fields, time.Unix(dat.Time, 0))
	if err != nil {
		return 0, err
	}
	out.buffer <- pt
	return len(p), nil
}

func (out *InfluxOutput) Close() error {
	return nil
}

func (out *InfluxOutput) Version() string {
	return "0.0.1"
}

func init() {
	transport.RegistOutputer("influxdb", NewInfluxOutput())
}
