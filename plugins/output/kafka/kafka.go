package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/luopengift/golibs/channel"
	"github.com/luopengift/golibs/logger"
	"github.com/luopengift/transport"
)

const (
	VERSION = "0.0.1"
)

type KafkaOutput struct {
	Addrs    []string `json:"addrs"`
	Topic    string   `json:"topic"`
	MaxProcs int      `json:"max_procs"` //最大并发写协程

	Message chan []byte //将数据写入这个管道中
	// 并发写topic的协程控制
	// 由于并发写入topic,写入顺序不可控,想要严格数序的话,maxThreads = 1即可
	channel *channel.Channel //并发写topic的协程控制
}

func NewKafkaOutput() *KafkaOutput {
	return new(KafkaOutput)
}

func (out *KafkaOutput) Init(config transport.Configer) error {
	err := config.Parse(out)
	if err != nil {
		return err
	}
	out.Message = make(chan []byte, out.MaxProcs)
	out.channel = channel.NewChannel(out.MaxProcs)
	return nil
}

func (out *KafkaOutput) ChanInfo() string {
	return out.channel.String()
}

func (out *KafkaOutput) Write(msg []byte) (int, error) {
	out.Message <- msg
	return len(msg), nil
}

func (out *KafkaOutput) Close() error {
	return out.Close()
}

func (out *KafkaOutput) Start() error {
	go out.WriteToTopic()
	return nil
}

func (out *KafkaOutput) WriteToTopic() error {

	config := sarama.NewConfig()
	config.ClientID = "TransportKafkaOutput"
	config.Producer.Return.Successes = true
	if err := config.Validate(); err != nil {
		logger.Error("<config error> %v", err)
		return err
	}

	producer, err := sarama.NewSyncProducer(out.Addrs, config)
	if err != nil {
		logger.Error("<Failed to produce message> %v", err)
		return err
	}
	defer producer.Close()

	for {
		select {
		case message := <-out.Message:
			out.channel.Add()
			go func(message []byte) {
				msg := &sarama.ProducerMessage{
					Topic: out.Topic,
					//Partition: int32(-1),
					//Key:       sarama.StringEncoder("key"),
					Value: sarama.ByteEncoder(message),
				}
				if partition, offset, err := producer.SendMessage(msg); err != nil {
					logger.Error("<write to kafka error,partition=%v,offset=%v> %v", partition, offset, err)
				}
				out.channel.Done()
			}(message)
		}
	}
}

func (out *KafkaOutput) Version() string {
	return VERSION
}

func init() {
	transport.RegistOutputer("kafka", NewKafkaOutput())
}
