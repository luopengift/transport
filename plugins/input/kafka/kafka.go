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
    *KafkaConfig
    Message chan []byte //从这个管道中读取数据
}

type KafkaConfig struct {
    Addrs   []string
    Topic   string
    Offset  int64
}



func NewKafkaInput() *KafkaInput {
    k := new(KafkaInput)
    return k
}

func (k *KafkaInput) Init(config pipeline.Configer) error {
    cfg := KafkaConfig{}
    err := config.Parse(&cfg)
    if err != nil {
        logger.Error("parse error:%v", err)
        return err
    }
    k.KafkaConfig = &cfg
    k.Message = make(chan []byte)
    return nil
}

func (k *KafkaInput) Read(p []byte) (cnt int, err error) {
    msg := <-k.Message
    n := copy(p, msg)
    return n, nil
}

func (k *KafkaInput) Start() error {
    k.ReadFromTopic()
    return nil
}

func (k *KafkaInput) ReadFromTopic() {
    var wg sync.WaitGroup
    consumer, err := sarama.NewConsumer(k.KafkaConfig.Addrs, sarama.NewConfig())
    if err != nil {
        logger.Error("<new consumer error> %v", err)
    }
    partitionList, err := consumer.Partitions(k.KafkaConfig.Topic)
    if err != nil {
        logger.Error("<consumer partitions> %v", err)
    }
    for partition := range partitionList {
        pc, err := consumer.ConsumePartition(k.KafkaConfig.Topic, int32(partition), k.KafkaConfig.Offset)
        if err != nil {
            logger.Error("<consume error> %v", err)
        }
        defer pc.AsyncClose()

        wg.Add(1)
        go func(pc sarama.PartitionConsumer) {
            defer wg.Done()
            for msg := range pc.Messages() {
                k.Message <- msg.Value
            }
        }(pc)

    }
    wg.Wait()
    consumer.Close()
}

func (k *KafkaInput) Close() error {
    return k.Close()
}

func init() {
    pipeline.RegistInputer("kafka", NewKafkaInput())
}
