package main

import (
	"fmt"
	"github.com/luopengift/transport"
	"github.com/luopengift/transport/plugins/kafka"
	//"os"
)

const (
	VERSION = "0.0.1"
)

func main() {
	//	file, _ := os.Open("/tmp/port.py")

	output := kafka.NewProducer([]string{"172.31.4.53:9092", "172.31.4.54:9092", "172.31.4.55:9092"}, "lp_test", 100)
	go output.WriteToTopic()

	input := kafka.NewConsumer([]string{"10.10.20.14:9092", "10.10.20.15:9092", "10.10.20.16:9092"}, "falcon_monitor_us", -1)
	go input.ReadFromTopic()

	//t := transport.NewTransport(input, os.Stdout)
	t := transport.NewTransport(input, output)
	t.Run()
	fmt.Printf("transport %s start success...\n", VERSION)
}
