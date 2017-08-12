package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/luopengift/golibs/channel"
	"github.com/luopengift/golibs/logger"
	"github.com/luopengift/transport/pipeline"
)

type KafkaOutput struct {
	*KafkaOutputConfig
	Message chan []byte //将数据写入这个管道中
	// 并发写topic的协程控制
	// 由于并发写入topic,写入顺序不可控,想要严格数序的话,maxThreads = 1即可
	channel *channel.Channel //并发写topic的协程控制
}

type KafkaOutputConfig struct {
	Addrs    []string `json:"addrs"`
	Topic    string   `json:"topic"`
	MaxBytes int      `json:"max_bytes"` //最大写入字节长度
	MaxProcs int      `json:"max_procs"` //最大并发写协程
}

func NewKafkaOutput() *KafkaOutput {
	return new(KafkaOutput)
}

func (k *KafkaOutput) Init(config pipeline.Configer) error {
	cfg := &KafkaOutputConfig{}
	err := config.Parse(cfg)
	if err != nil {
		return err
	}
	k.KafkaOutputConfig = cfg
	k.Message = make(chan []byte, k.MaxBytes)
	k.channel = channel.NewChannel(k.MaxProcs)
	return nil
}

func (k *KafkaOutput) ChanInfo() string {
	return k.channel.String()
}

func (k *KafkaOutput) Write(msg []byte) (int, error) {
	k.Message <- msg
	return len(msg), nil
}

func (k *KafkaOutput) Close() error {
	return k.Close()
}

func (k *KafkaOutput) Start() error {
	go k.WriteToTopic()
	return nil
}

func (k *KafkaOutput) WriteToTopic() error {

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	if err := config.Validate(); err != nil {
		logger.Error("<config error> %v", err)
		return err
	}

	producer, err := sarama.NewSyncProducer(k.KafkaOutputConfig.Addrs, config)
	if err != nil {
		logger.Error("<Failed to produce message> %v", err)
		return err
	}
	defer producer.Close()

	for {
		select {
		case message := <-k.Message:
			k.channel.Add()
			go func(message []byte) {
				msg := &sarama.ProducerMessage{
					Topic:     k.KafkaOutputConfig.Topic,
					Partition: int32(-1),
					Key:       sarama.StringEncoder("key"),
					Value:     sarama.ByteEncoder(message),
				}
				if partition, offset, err := producer.SendMessage(msg); err != nil {
					logger.Error("<write to kafka error,partition=%v,offset=%v> %v", partition, offset, err)
				}
				k.channel.Done()
			}(message)
		}
	}
	return nil
}

func init() {
	pipeline.RegistOutputer("kafka", NewKafkaOutput())
}
