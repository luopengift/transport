package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/luopengift/golibs/channel"
	"github.com/luopengift/golibs/logger"
	"github.com/luopengift/transport"
	"strings"
)

type KafkaOutput struct {
	addrs   []string
	topic   string
	message chan []byte //将数据写入这个管道中
	// 并发写topic的协程控制
	// 由于并发写入topic,写入顺序不可控,想要严格数序的话,maxThreads = 1即可
	channel *channel.Channel //并发写topic的协程控制
}

func NewKafkaOutput() *KafkaOutput {
	return new(KafkaOutput)
}

func (k *KafkaOutput) Init(cfg map[string]string) error {
	k.addrs = strings.Split(cfg["addrs"], ",")
	k.topic = cfg["topic"]
	k.message = make(chan []byte, 10000)
	k.channel = channel.NewChannel(1000000)
	return nil
}

func (self *KafkaOutput) ChanInfo() string {
	return self.channel.String()
}

func (self *KafkaOutput) Write(msg []byte) (int, error) {
	self.message <- msg
	return len(msg), nil
}

func (self *KafkaOutput) Close() error {
	return self.Close()
}

func (self *KafkaOutput) Start() error {
	go self.WriteToTopic()
	return nil
}

func (self *KafkaOutput) WriteToTopic() error {

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	if err := config.Validate(); err != nil {
		logger.Error("<config error> %v", err)
		return err
	}

	producer, err := sarama.NewSyncProducer(self.addrs, config)
	if err != nil {
		logger.Error("<Failed to produce message> %v", err)
		return err
	}
	defer producer.Close()

	for {
		select {
		case message := <-self.message:
			self.channel.Add()
			go func(message []byte) {
				msg := &sarama.ProducerMessage{
					Topic:     self.topic,
					Partition: int32(-1),
					Key:       sarama.StringEncoder("key"),
					Value:     sarama.ByteEncoder(message),
				}
				if partition, offset, err := producer.SendMessage(msg); err != nil {
					logger.Error("<write to kafka error,partition=%v,offset=%v> %v", partition, offset, err)
				}
				self.channel.Done()
			}(message)
		}
	}
	return nil
}

func init() {
	transport.RegistOutputer("kafka", NewKafkaOutput())
}
