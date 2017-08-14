package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/luopengift/golibs/logger"
	"github.com/luopengift/transport/pipeline"
	"sync"
)

/*
一些使用说明:
sarame.OffsetNewest int64 = -1
sarame.OffsetOldest int64 = -2
*/
type KafkaInput struct {
	Addrs  []string
	Topic  string
	Offset int64

    Message chan []byte //从这个管道中读取数据
}


func NewKafkaInput() *KafkaInput {
	in := new(KafkaInput)
	return in
}

func (in *KafkaInput) Init(config pipeline.Configer) error {
	err := config.Parse(in)
	if err != nil {
		logger.Error("parse error:%v", err)
		return err
	}
	in.Message = make(chan []byte)
	return nil
}

func (in *KafkaInput) Read(p []byte) (cnt int, err error) {
	msg := <-in.Message
	n := copy(p, msg)
	return n, nil
}

func (in *KafkaInput) Start() error {
	in.ReadFromTopic()
	return nil
}

func (in *KafkaInput) ReadFromTopic() {
	var wg sync.WaitGroup
	consumer, err := sarama.NewConsumer(in.Addrs, sarama.NewConfig())
	if err != nil {
		logger.Error("<new consumer error> %v", err)
	}
	partitionList, err := consumer.Partitions(in.Topic)
	if err != nil {
		logger.Error("<consumer partitions> %v", err)
	}
	for partition := range partitionList {
		pc, err := consumer.ConsumePartition(in.Topic, int32(partition), in.Offset)
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

func (in *KafkaInput) Close() error {
	return in.Close()
}

func init() {
	pipeline.RegistInputer("kafka", NewKafkaInput())
}
