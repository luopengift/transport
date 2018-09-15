package kafka

import (
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/luopengift/golibs/logger"
	"github.com/luopengift/transport"
	"github.com/wvanbergen/kafka/consumergroup"
	"github.com/wvanbergen/kazoo-go"
)

const (
	VERSION = "0.0.1"
)

// KafkaInput kafka input
// sarame.OffsetNewest int64 = -1
// sarame.OffsetOldest int64 = -2
type KafkaInput struct {
	Addrs  []string `json:"addrs" yaml:"addr"` //如果定义了group,则addrs是zookeeper的地址(2181)，否则的话是kafka的地址(9092)
	Topics []string `json:"topics" yaml:"topics"`
	Group  string   `json:"group" yaml:"group"`
	Offset int64    `json:"offset" yaml:"offset"`

	Message chan []byte //从这个管道中读取数据
}

func NewKafkaInput() *KafkaInput {
	in := new(KafkaInput)
	return in
}

func (in *KafkaInput) Init(config transport.Configer) error {
	err := config.Parse(in)
	if err != nil {
		logger.Error("parse error:%v", err)
		return err
	}
	in.Message = make(chan []byte, 1000)
	return nil
}

func (in *KafkaInput) Read(p []byte) (cnt int, err error) {
	msg := <-in.Message
	n := copy(p, msg)
	return n, nil
}

func (in *KafkaInput) Start() error {
	if in.Group == "" {
		for _, topic := range in.Topics {
			go in.ReadFromTopic(topic)
		}
	} else {
		in.ReadWithGroup()
	}
	return nil
}

// 简单kafka消费者
func (in *KafkaInput) ReadFromTopic(topic string) {
	var wg sync.WaitGroup
	consumer, err := sarama.NewConsumer(in.Addrs, sarama.NewConfig())
	if err != nil {
		logger.Error("<new consumer error> %v", err)
	}
	partitionList, err := consumer.Partitions(topic)
	if err != nil {
		logger.Error("<consumer partitions> %v", err)
	}
	for partition := range partitionList {
		pc, err := consumer.ConsumePartition(topic, int32(partition), in.Offset)
		if err != nil {
			logger.Error("<consume error> %v", err)
		}
		defer pc.AsyncClose()

		wg.Add(1)
		go func(pc sarama.PartitionConsumer) {
			defer wg.Done()
			for msg := range pc.Messages() {
				in.Message <- msg.Value
			}
		}(pc)

	}
	wg.Wait()
	consumer.Close()
}

// 多个consumer group
func (in *KafkaInput) ReadWithGroup() error {
	config := consumergroup.NewConfig()
	config.Offsets.Initial = in.Offset
	config.Offsets.ProcessingTimeout = 10 * time.Second

	zookeeperNodes, chroot := kazoo.ParseConnectionString(in.Addrs[0])
	config.Zookeeper.Chroot = chroot
	consumer, err := consumergroup.JoinConsumerGroup(in.Group, in.Topics, zookeeperNodes, config)
	if err != nil {
		logger.Error("parse error:%v", err)
		return err
	}

	go func() {
		for err := range consumer.Errors() {
			logger.Error("consumer error:%v", err)
		}
	}()

	eventCount := 0
	offsets := make(map[string]map[int32]int64)

	for message := range consumer.Messages() {
		if offsets[message.Topic] == nil {
			offsets[message.Topic] = make(map[int32]int64)
		}

		eventCount += 1
		if offsets[message.Topic][message.Partition] != 0 && offsets[message.Topic][message.Partition] != message.Offset-1 {
			logger.Error("Unexpected offset on %s:%d. Expected %d, found %d, diff %d.", message.Topic, message.Partition, offsets[message.Topic][message.Partition]+1, message.Offset, message.Offset-offsets[message.Topic][message.Partition]+1)
		}

		offsets[message.Topic][message.Partition] = message.Offset
		consumer.CommitUpto(message)
		in.Message <- message.Value
	}

	logger.Info("Processed %d events.", eventCount)
	logger.Info("%+v", offsets)
	consumer.Close()
	return nil

}

func (in *KafkaInput) Close() error {
	return in.Close()
}

func (in *KafkaInput) Version() string {
	return VERSION
}

func init() {
	transport.RegistInputer("kafka", NewKafkaInput())
}
