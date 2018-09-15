package transport

import (
	"fmt"
	"os"
	"strings"
	"time"

	//"github.com/luopengift/golibs/logger"
	"github.com/luopengift/log"
)

// Transport core struct
type Transport struct {
	byteSize int
	chanSize int

	Inputs   []*Input
	Outputs  []*Output
	Adapts   []*Adapt
	recvChan chan []byte
	sendChan chan []byte
	isEnd    chan bool
	logs     *log.Log
}

// NewTransprot new transport
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
	transport.logs = log.NewLog("transport", os.Stdout)
	transport.logs.SetTimeFormat("2006/01/02 15:03:04.000")
	if !cfg.Runtime.DEBUG {
		transport.logs.SetLevel(log.INFO)
	}
	startCronTask()

	return transport, err
}

func (t *Transport) injectInput(p []byte) {
	t.recvChan <- p
}

func (t *Transport) injectOutput(p []byte) {
	t.sendChan <- p
}

// RunInputs run input plugins
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

// RunAdapts run adapt plugins
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
					t.logs.Error("[%s] %s", c.Name, ErrReadBufferClosedError.Error())
					break
				}
			}
		}(adapt)
	}
}

// RunOutputs run output plugins
func (t *Transport) RunOutputs() {
	for _, output := range t.Outputs {
		go output.Start()
	}
	for {
		value, ok := <-t.sendChan
		if !ok {
			t.logs.Error("%s", ErrWriteBufferClosedError.Error())
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

// Run run
func (t *Transport) Run() {
	go t.RunOutputs()
	go t.RunAdapts()
	go t.RunInputs()
}

// Stop stop
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

// Stat monitor
func (t *Transport) Stat() string {
	inputStat := []string{}
	for _, input := range t.Inputs {
		inputStat = append(inputStat, fmt.Sprintf("%v:%v", input.Name, input.Count()))
	}
	adaptStat := []string{}
	for _, adapt := range t.Adapts {
		adaptStat = append(adaptStat, fmt.Sprintf("%v:%v", adapt.Name, adapt.Count()))
	}
	outputStat := []string{}
	for _, output := range t.Outputs {
		outputStat = append(outputStat, fmt.Sprintf("%v:%v", output.Name, output.Count()))
	}
	return fmt.Sprintf("stat=> [inputs]:%s|[adapts]:%s|[outputs]:%s", strings.Join(inputStat, ","), strings.Join(adaptStat, ","), strings.Join(outputStat, ","))
}

// T transport
var T *Transport
