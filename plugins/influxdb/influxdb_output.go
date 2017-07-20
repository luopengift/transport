package influxdb

import (
	"encoding/json"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/luopengift/transport"
	"time"
)

type Data struct {
	Timestamp time.Time              `json:"time"`
	Name      string                 `json:"name"`
	Fields    map[string]interface{} `json:"fields"`
	Tags      map[string]string      `json:"tags"`
}

type InfluxOutput struct {
	client.Client
	client.BatchPoints
}

func NewInfluxOutput() *InfluxOutput {
	influx := new(InfluxOutput)
	return influx
}

func (influx *InfluxOutput) Init(cfg map[string]string) error {
	var err error
	influx.Client, err = client.NewHTTPClient(client.HTTPConfig{
		Addr: cfg["addr"],
		//Addr: "http://localhost:8086",
		Username:           cfg["username"],
		Password:           cfg["password"],
		UserAgent:          "InfluxDBClient",
		Timeout:            5 * time.Second,
		InsecureSkipVerify: false,
	})

	if err != nil {
		return err
	}

	// Create a new point batch
	influx.BatchPoints, err = client.NewBatchPoints(client.BatchPointsConfig{
		Database: cfg["database"],
		//Precision is the write precision of the points, defaults to "ns"
		Precision: "s",
	})

	if err != nil {
		return err
	}

	return nil
}

func (influx *InfluxOutput) Start() error {
	return nil
}
func (influx *InfluxOutput) Write(p []byte) (int, error) {
	data := Data{}
	err := json.Unmarshal(p, &data)
	if err != nil {
		return 0, err
	}
	// Create a point and add to batch
	pt, err := client.NewPoint(data.Name, data.Tags, data.Fields, data.Timestamp)

	if err != nil {
		return 0, err
	}

	influx.BatchPoints.AddPoint(pt)

	// Write the batch
	err = influx.Client.Write(influx.BatchPoints)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (influx *InfluxOutput) Close() error {
	return influx.Client.Close()
}

func init() {
	transport.RegistOutputer("influxdb", NewInfluxOutput())
}
