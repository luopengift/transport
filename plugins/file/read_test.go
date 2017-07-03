package file

import (
	"fmt"
	"testing"
)

func Test_read(t *testing.T) {
	//tt := NewTail("test-%Y-%M-%D-%h-%m.log")
	tt := NewTail("test.log")
	tt.ReadLine()
	tt.EndStop(true)
	for v := range tt.NextLine() {
		fmt.Println(*v)
	}
}

type RuntimeConfig struct {
	DEBUG    bool `json:"DEBUG"`
	MAXPROCS int  `json:"MAXPROCS"`
}

type KafkaConfig struct {
	Addrs      []string `json:"addrs"`
	Topic      string   `json:"topic"`
	MaxThreads int64    `json:"maxthreads"`
}

type HttpConfig struct {
	Addr string `json:"addr"`
}

type TestConfig struct {
	Runtime RuntimeConfig
	Kafka   KafkaConfig
	File    []string `json:"file"`
	Prefix  string   `json:"prefix"`
	Suffix  string   `json:"suffix"`
	Http    HttpConfig
	Tags    string `json:"tags"`
	Version string `json:version`
}

func Test_config(t *testing.T) {
	test := &TestConfig{}
	config := NewConfig("./config.json")
	config.Parse(test)
	fmt.Println(fmt.Sprintf("%+v", test))
	fmt.Println(config)

}
