package influxdb

import (
	"fmt"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/influxdata/influxdb/models"
	"github.com/luopengift/log"
	"github.com/luopengift/transport"
	"github.com/luopengift/types"
)

type InfluxInput struct {
	Addr        string `json:"addr"`
	DB          string `json:"database"`
	Precision   string `json:"precision"`
	User        string `json:"username"`
	Pass        string `json:"password"`
	QueryString string `json:"query"`
	Buffer      chan []interface{}
	client      client.Client
}

func NewInfluxInput() *InfluxInput {
	return new(InfluxInput)
}

func (in *InfluxInput) Query(str string) (models.Row, error) {
	response, err := in.client.Query(client.Query{
		Command:  str,
		Database: in.DB,
	})
	if err != nil {
		return models.Row{}, fmt.Errorf("influxdb query error:%v", err)
	}
	if response.Error() != nil {
		return models.Row{}, fmt.Errorf("influxdb response error:%v", response.Error())
	}
	result := response.Results[0]
	if result.Err != "" {
		return models.Row{}, fmt.Errorf("influxdb result error:%v", result.Err)
	}
	return result.Series[0], nil
}

func (in *InfluxInput) Init(cfg transport.Configer) error {
	in.DB = "mydb"
	in.Precision = "ns"
	err := cfg.Parse(in)
	if err != nil {
		return err
	}

	in.client, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     in.Addr,
		Username: in.User,
		Password: in.Pass,
	})
	in.Buffer = make(chan []interface{}, 1000)
	return err
}

func (in *InfluxInput) Start() error {
	data, err := in.Query(in.QueryString)
	if err != nil {
		return err
	}
	log.Debug("%#v, %#v", data.Name, data.Columns)
	for _, dat := range data.Values {
		in.Buffer <- dat
	}
	return nil
}

func (in *InfluxInput) Read(p []byte) (int, error) {
	b, err := types.ToBytes(<-in.Buffer)
	if err != nil {
		return 0, err
	}
	n := copy(p, b)
	return n, nil
}

func (in *InfluxInput) Close() error {
	return in.Close()
}

func (in *InfluxInput) Version() string {
	return "0.0.1"
}

func init() {
	transport.RegistInputer("influxdb", NewInfluxInput())
}
