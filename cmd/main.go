package main

import (
	"github.com/luopengift/golibs/logger"
	"github.com/luopengift/transport"
	"github.com/luopengift/transport/plugins/hdfs"
	"github.com/luopengift/transport/plugins/kafka"
	"github.com/luopengift/transport/filter"
	"time"
)

const (
	VERSION = "0.0.1"
)

var output *hdfs.HDFS
var t *transport.Transport

func main() {
	//	file, _ := os.Open("/tmp/port.py")
    logger.Info("Transport starting...")
	//output := kafka.NewProducer([]string{"172.31.4.53:9092", "172.31.4.54:9092", "172.31.4.55:9092"}, "lp_test", 100)
	//go output.WriteToTopic()

	//input := kafka.NewConsumer([]string{"172.31.4.53:9092", "172.31.4.54:9092", "172.31.4.55:9092"}, "falcon_monitor_ap", -1)
	input := kafka.NewConsumer([]string{"172.31.4.53:9092", "172.31.4.54:9092", "172.31.4.55:9092"}, "lp_test", -1)
	//input := kafka.NewConsumer([]string{"10.10.20.14:9092", "10.10.20.15:9092", "10.10.20.16:9092"}, "lp_test", -1)
	//input := kafka.NewConsumer([]string{"10.10.20.14:9092", "10.10.20.15:9092", "10.10.20.16:9092"}, "falcon_monitor_us", -1)
	go input.ReadFromTopic()

	output = hdfs.NewHDFS("10.10.20.64:8020", "/tmp/luopeng/%Y%M%D/%h", "test.log")
	err := output.Init()
    if err != nil {
        logger.Error("Init error:%v",err)
        return
    }	
    defer output.Close()
	//t := transport.NewTransport(input, os.Stdout)
	t = transport.NewTransport(input, filter.AddEnter, output)
    t.StartWrite()
	go func() {
        tc := time.Tick(10*time.Second)
		for {
			select {
			case <-tc:
                t.StopWrite()
                stime := time.Now()
                //time.Sleep(1 * time.Second)
				logger.Info("prepare start:%#v,%v",output,stime)
				if err := output.Close(); err != nil {
                    logger.Error("!close error:%v",err)
                }
                
	            output = hdfs.NewHDFS("10.10.20.64:8020", "/tmp/luopeng/%Y%M%D/%h", "test.log")
                output.Init()
				t.SetOutputer(output)
                logger.Info("prepare end:%#v,%v",output,time.Since(stime))
                t.StartWrite()
			}
		}
	}()
	t.Run()
}