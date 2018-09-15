package kafka2

import (
	"context"
	"fmt"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/luopengift/log"
	"github.com/luopengift/transport"
)

const (
	// Version version
	VERSION = "0.0.1"
)

// KafkaInput
// 一些使用说明:
// sarame.OffsetNewest int64 = -1
// sarame.OffsetOldest int64 = -2
type KafkaInput struct {
	Addrs   []string    `json:"addrs" yaml:"addrs"` //如果定义了group,则addrs是zookeeper的地址(2181)，否则的话是kafka的地址(9092)
	Topics  []string    `json:"topics" yaml:"topics"`
	Group   string      `json:"group" yaml:"group"`
	Offset  string      `json:"offset" yaml:"offset"`
	Message chan []byte //从这个管道中读取数据
	*kafka.Consumer
}

// NewKafkaInput kafka input
func NewKafkaInput() *KafkaInput {
	in := new(KafkaInput)
	return in
}

func (in *KafkaInput) Init(config transport.Configer) error {
	err := config.Parse(in)
	if err != nil {
		log.Error("parse error:%v", err)
		return err
	}
	in.Message = make(chan []byte, 1000)

	//c.closechan = make(chan string, 1)
	//c.mux = new(sync.Mutex)
	if in.Consumer, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": strings.Join(in.Addrs, ","),
		"group.id":          in.Group,
		"auto.offset.reset": in.Offset,
	}); err != nil {
		return err
	}
	in.Consumer.SubscribeTopics(in.Topics, nil)
	return nil

}

func (in *KafkaInput) Read(p []byte) (cnt int, err error) {
	msg := <-in.Message
	n := copy(p, msg)
	return n, nil
}

func (in *KafkaInput) ReadFromTopics(ctx context.Context) error {
LOOP:
	for {
		select {
		case <-ctx.Done():
			break LOOP
		default:
			msg, err := in.Consumer.ReadMessage(-1)
			if err != nil {
				fmt.Printf("Consumer error: %v (%v)\n", err, msg)
				break LOOP
			}
			in.Message <- msg.Value
			//fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value)[:60])
		}
	}
	log.Warn("close client!")
	close(in.Message)
	return in.Consumer.Close()
}

func (in *KafkaInput) Start() error {
	go in.ReadFromTopics(context.TODO())
	return nil
}

func (in *KafkaInput) Close() error {
	return in.Close()
}

func (in *KafkaInput) Version() string {
	return VERSION
}

func init() {
	transport.RegistInputer("kafka2", NewKafkaInput())
}
