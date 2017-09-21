package transport

import (
	"fmt"
	"github.com/luopengift/golibs/logger"
	"os"
	"strings"
	"time"
)

type Transport struct {
	byteSize int
	chanSize int

	Inputs   []*Input
	Outputs  []*Output
	Adapts   []*Adapt
	recvChan chan []byte
	sendChan chan []byte
	isEnd    chan bool
	logs     *logger.Logger
}

func NewTransport(cfg *Config) (*Transport, error) {
	var err error
	transport := new(Transport)
	transport.Inputs, err = cfg.InitInputs()
	if err != nil {
		return nil, err
	}
	transport.Adapts, err = cfg.InitAdapts()
	if err != nil {
		return nil, err
	}
	transport.Outputs, err = cfg.InitOutputs()
	if err != nil {
		return nil, err
	}
	transport.byteSize = cfg.Runtime.BYTESIZE
	transport.chanSize = cfg.Runtime.CHANSIZE
	transport.recvChan = make(chan []byte, transport.chanSize)
	transport.sendChan = make(chan []byte, transport.chanSize)
	transport.isEnd = make(chan bool)
	if cfg.Runtime.DEBUG {
		transport.logs = logger.NewLogger(logger.DEBUG, os.Stdout)
	} else {
		transport.logs = logger.NewLogger(logger.INFO, os.Stdout)
	}
	transport.logs.SetPrefix("[transport]")

	startCronTask()

	return transport, err
}

func (t *Transport) injectInput(p []byte) {
	t.recvChan <- p
}

func (t *Transport) injectOutput(p []byte) {
	t.sendChan <- p
}

func (t *Transport) RunInputs() {
	for _, input := range t.Inputs {
		go input.Inputer.Start()
		go func(in *Input) {
			for {
				b := make([]byte, t.byteSize, t.byteSize)
				n, err := in.Read(b)
				if err != nil {
					t.logs.Error("[%s] read error:%s", in.Name, err.Error())
					continue
				}
				t.recvChan <- b[:n]
				t.logs.Debug("[%s] recv %v", in.Name, string(b[:n]))
			}
		}(input)
	}

}

func (t *Transport) RunAdapts() {
	for _, adapt := range t.Adapts {
		go func(c *Adapt) {
			for {
				if value, ok := <-t.recvChan; ok {
					c.channel.Add()
					go func() {
						b := make([]byte, t.byteSize, t.byteSize)
						n, err := c.Handle(value, b)
						if err != nil {
							t.logs.Error("[%s] %s", c.Name, err.Error())
						} else {
							t.sendChan <- b[:n]
						}
						c.channel.Done()
					}()
				} else {
					t.logs.Error("[%s] %s", c.Name, ReadBufferClosedError.Error())
					break
				}
			}
		}(adapt)
	}
}

func (t *Transport) RunOutputs() {
	for _, output := range t.Outputs {
		go output.Start()
	}
	for {
		value, ok := <-t.sendChan
		if !ok {
			t.logs.Error("%s", WriteBufferClosedError.Error())
			break
		}
		t.logs.Debug("send %v", string(value))
		for _, output := range t.Outputs {
			func(out *Output) {
				n, err := out.Write(value)
				if err != nil {
					t.logs.Error("[%s] write data err:%s,%v", out.Name, err.Error(), n)
				}
			}(output)
		}
	}
}

func (t *Transport) Run() {
	go t.RunOutputs()
	go t.RunAdapts()
	go t.RunInputs()
}

func (t *Transport) Stop() {
	stopCronTask() //关闭全部定时任务
	for _, input := range t.Inputs {
		input.Inputer.Close()
	}
	close(t.recvChan)
	close(t.sendChan)
	for _, output := range t.Outputs {
		output.Outputer.Close()
	}
	t.logs.Info("stop success...%s", time.Now())
}

func (t *Transport) Stat() string {
	input_stat := []string{}
	for _, input := range t.Inputs {
		input_stat = append(input_stat, fmt.Sprintf("%v:%v", input.Name, input.Count()))
	}
	adapt_stat := []string{}
	for _, adapt := range t.Adapts {
		adapt_stat = append(adapt_stat, fmt.Sprintf("%v:%v", adapt.Name, adapt.Count()))
	}
	output_stat := []string{}
	for _, output := range t.Outputs {
		output_stat = append(output_stat, fmt.Sprintf("%v:%v", output.Name, output.Count()))
	}
	return fmt.Sprintf("stat=> [inputs]:%s|[adapts]:%s|[outputs]:%s", strings.Join(input_stat, ","), strings.Join(adapt_stat, ","), strings.Join(output_stat, ","))
}

var T *Transport
