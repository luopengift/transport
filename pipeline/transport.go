package pipeline

import (
	"fmt"
	"github.com/luopengift/golibs/logger"
	"os"
	"time"
)

type Transport struct {
	Inputs    []*Input
	Outputs   []*Output
	Codecs    []*Codec
	recv_chan chan []byte
	send_chan chan []byte
	errchan   chan string
	isEnd     chan bool
	logs      *logger.Logger
}

func NewTransport(cfg *Config) *Transport {
	transport := new(Transport)
	transport.Inputs = cfg.InitInputs()
	transport.Codecs = cfg.InitCodecs()
	transport.Outputs = cfg.InitOutputs()
	transport.recv_chan = make(chan []byte, 100)
	transport.send_chan = make(chan []byte, 100)
	transport.errchan = make(chan string, 100)
	transport.isEnd = make(chan bool)
	transport.logs = logger.NewLogger(logger.DEBUG, os.Stdout)
	transport.logs.SetPrefix("[transport]")

	return transport
}

func (t *Transport) RunInputs() {
	for _, input := range t.Inputs {
		go input.Inputer.Start()
		go func(in *Input) {
			for {
				//	in.Mutex.Lock()
				//	defer in.Mutex.Unlock()
				b := make([]byte, MAX, MAX)
				n, err := in.Read(b)
				if err != nil {
					t.errchan <- fmt.Sprintf("[%s] %s", in.Name, err.Error())
					continue
				}
				t.recv_chan <- b[:n]
				t.logs.Debug("recv %v", string(b[:n]))
			}
		}(input)
	}

}

func (t *Transport) RunCodecs() {
	for _, codec := range t.Codecs {
		go func(h *Codec) {
			for {
				value, ok := <-t.recv_chan
				if !ok {
					t.errchan <- fmt.Sprintf("[%s] %s", h.Name, ReadBufferClosedError.Error())
					t.logs.Error("[%s] %s", h.Name, ReadBufferClosedError.Error())
					break
				}
				b := make([]byte, MAX, MAX)
				n, err := h.Handle(value, b)
				if err != nil {
					t.errchan <- fmt.Sprintf("[%s] %s", h.Name, err.Error())
					continue
				}
				t.send_chan <- b[:n]
			}
		}(codec)
	}
}

func (t *Transport) RunOutputs() {
	for _, output := range t.Outputs {
		go output.Start()
	}
	for {
		value, ok := <-t.send_chan
		if !ok {
			t.logs.Error("output send err:%v", ok)
			break
		}
		for _, output := range t.Outputs {
			func(out *Output) {
				_, err := out.Outputer.Write(value)
				if err != nil {
					t.errchan <- fmt.Sprintf("[%s] write data err:%s", out.Name, err.Error())
				}
			}(output)
		}
	}
}

func (t *Transport) Run() {
	go t.RunInputs()
	go t.RunOutputs()
	go t.RunCodecs()
}

func (t *Transport) Stop() {
	for _,input := range t.Inputs {
        input.Inputer.Close()
    }
	close(t.recv_chan)
	close(t.send_chan)
    for _,output := range t.Outputs {
        output.Outputer.Close()
    }
	t.logs.Info("stop success...%s", time.Now())
}
