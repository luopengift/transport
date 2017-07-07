package kafka

import (
	"errors"
	"github.com/Shopify/sarama"
	"github.com/luopengift/golibs/logger"
	"sync"
)

/*
一些使用说明:
sarame.OffsetNewest int64 = -1
sarame.OffsetOldest int64 = -2
*/
type Consumer struct {
	Addrs   []string
	Topic   string
	Offset  int64
	Message chan *[]byte //从这个管道中读取数据
}

func NewConsumer(addrs []string, topic string, offset int64) *Consumer {
	return &Consumer{
		Addrs:   addrs,
		Topic:   topic,
		Offset:  offset,
		Message: make(chan *[]byte),
	}
}

func (self *Consumer) Read(p []byte) (cnt int, err error) {
	msg := <-self.Message
	if len(*msg) > len(p) {
		p = (*msg)[:len(p)-1]
		return len(p), errors.New("message is larger than buffer")
	}
	copy(p, *msg)
	return len(*msg), nil
}

func (self *Consumer) Start() error {
    self.ReadFromTopic()
    return nil
}


func (self *Consumer) ReadFromTopic() {
	var wg sync.WaitGroup
	//consumer
	consumer, err := sarama.NewConsumer(self.Addrs, sarama.NewConfig())
	if err != nil {
		logger.Error("<new consumer error> %v", err)
	}
	//topics,err := consumer.Topics() //获取topic列表
	partitionList, err := consumer.Partitions(self.Topic)
	if err != nil {
		logger.Error("<consumer partitions> %v", err)
	}
	for partition := range partitionList {
		pc, err := consumer.ConsumePartition(self.Topic, int32(partition), self.Offset)
		if err != nil {
			logger.Error("<consume error> %v", err)
		}
		defer pc.AsyncClose()

		wg.Add(1)

		go func(pc sarama.PartitionConsumer) {
			defer wg.Done()
			for msg := range pc.Messages() {
				self.Message <- &(msg.Value)
				//continue
				//fmt.Println("Partition:%d, Offset:%d, Key:%s, Value:%s", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
			}
		}(pc)

	}
	wg.Wait()
	//logs.Info("Done consuming topic", self.Topic)
	consumer.Close()
}

func (self *Consumer) Close() error {
	return self.Close()
}
