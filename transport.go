package transport

import (
	"fmt"
	"github.com/luopengift/golibs/logger"
	"os"
	"strings"
	"time"
)

type Transport struct {
	Inputs    []*Input
	Outputs   []*Output
	Codecs    []*Codec
	recv_chan chan []byte
	send_chan chan []byte
	isEnd     chan bool
	logs      *logger.Logger
}

func NewTransport(cfg *Config) (*Transport, error) {
    var err error
	transport := new(Transport)
	transport.Inputs, err = cfg.InitInputs()
	if err != nil {
        return nil, err
    }
    transport.Codecs, err = cfg.InitCodecs()
	if err != nil {
        return nil, err
    }
	transport.Outputs, err = cfg.InitOutputs()
	if err != nil {
        return nil, err
    }
	transport.recv_chan = make(chan []byte, 100)
	transport.send_chan = make(chan []byte, 100)
	transport.isEnd = make(chan bool)
	transport.logs = logger.NewLogger(logger.DEBUG, os.Stdout)
	transport.logs.SetPrefix("[transport]")

	startCronTask()

	return transport, err
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
					t.logs.Error("[%s] %s", in.Name, err.Error())
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
					t.logs.Error("[%s] %s", h.Name, ReadBufferClosedError.Error())
					break
				}
				h.Channel.Add()
				go func(value []byte) {
					defer h.Channel.Done()
					b := make([]byte, MAX, MAX)
					n, err := h.Handle(value, b)
					if err != nil {
						t.logs.Error("[%s] %s", h.Name, err.Error())
						return
					}
					t.send_chan <- b[:n]
				}(value)
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
				_, err := out.Write(value)
				if err != nil {
					t.logs.Error("[%s] write data err:%s", out.Name, err.Error())
				}
			}(output)
		}
	}
}

func (t *Transport) Run() {
	go func() {
		for {
			input_stat := []string{}
			for _, input := range t.Inputs {
				input_stat = append(input_stat, fmt.Sprintf("%v:%v", input.Name, input.Cnt))
			}
			codec_stat := []string{}
			for _, codec := range t.Codecs {
				codec_stat = append(codec_stat, fmt.Sprintf("%v:%v", codec.Name, codec.Cnt))
			}
			output_stat := []string{}
			for _, output := range t.Outputs {
				output_stat = append(output_stat, fmt.Sprintf("%v:%v", output.Name, output.Cnt))
			}
			t.logs.Info("stat=> inputs:%s|codecs:%s|outputs:%s", strings.Join(input_stat, ","), strings.Join(codec_stat, ","), strings.Join(output_stat, ","))
			time.Sleep(10 * time.Second)
		}
	}()
	go t.RunInputs()
	go t.RunCodecs()
	go t.RunOutputs()
}

func (t *Transport) Stop() {
	stopCronTask() //关闭全部定时任务
	for _, input := range t.Inputs {
		input.Inputer.Close()
	}
	close(t.recv_chan)
	close(t.send_chan)
	for _, output := range t.Outputs {
		output.Outputer.Close()
	}
	t.logs.Info("stop success...%s", time.Now())
}