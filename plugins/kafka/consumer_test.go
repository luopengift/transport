package kafka

import (
	"easyWork/framework/logs"
	"testing"
)

func Test_logs(t *testing.T) {
	c := NewConsumer([]string{"172.31.4.53:9092"}, "zhizi-monitor-data4", -1)
	go func() {
		for msg := range c.Message {
			logs.Info(string(*msg))
		}
	}()
	c.ReadFromTopic()
}
