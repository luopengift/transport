package pipeline

import (
	"fmt"
	"github.com/luopengift/golibs/logger"
	"os"
	"time"
)

type Transport struct {
	*Input
	*Output
	*Middleware
	r_chan chan []byte
	w_chan chan []byte
	isEnd  chan bool
	logs   *logger.Logger
}

func NewTransport(cfg *Config) *Transport {
	input := cfg.Input()
	output := cfg.Output()
	handle := cfg.Handle()

	transport := new(Transport)
	transport.Input = NewInput(input)
	transport.Middleware = NewMiddleware(handle, 2000)
	transport.Output = NewOutput(output)
	transport.r_chan = make(chan []byte, 1000)
	transport.w_chan = make(chan []byte, 1000)
	transport.isEnd = make(chan bool)
	transport.logs = logger.NewLogger(logger.INFO, "2006/01/02 15:04:05.000 [transport]", true, os.Stdout)

	return transport
}

func (t *Transport) End() chan bool {
	return t.isEnd
}

// 将数据从read接口读入 r_chan中
func (t *Transport) recv() error {
	t.Input.Mutex.Lock()
	defer t.Input.Mutex.Unlock()
	if t.Inputer == nil {
		return InputNullError
	}
	b := make([]byte, MAX, MAX)
	n, err := t.Inputer.Read(b)
	if err != nil {
		return err
	}
	t.r_chan <- b[:n]
	t.logs.Debug("recv %v", string(b[:n]))
	return nil
}

func (t *Transport) handle() error {
	tmp, ok := <-t.r_chan
	if !ok {
		return ReadBufferClosedError
	}
	t.Middleware.Channel.Add()
	go func(p []byte) {
		defer t.Middleware.Channel.Done()
		b := make([]byte, MAX, MAX)
		n, err := t.Handler.Handle(p, b)
		if err != nil {
			t.logs.Error("Handler Error!%v", err)
			return
		}
		t.w_chan <- b[:n]
	}(tmp)
	return nil
}

// 将数据从w_chan写入 Write接口中
func (t *Transport) send() error {
	t.Output.Mutex.Lock()
	defer t.Output.Mutex.Unlock()
	if t.Outputer == nil {
		return OutputNullError
	}
	write, ok := <-t.w_chan
	if !ok {
		return fmt.Errorf("w_chan is closed")
	}
	go func(b []byte) error {
		n, err := t.Outputer.Write(b)
		if err != nil {
			return fmt.Errorf("send error:%v,%v|message:%s", n, err, string(b))
		}
		t.logs.Debug("send:%v|msg:%v", n, string(b))
		return nil
	}(write)
	return nil
}

func (t *Transport) Run() {
	go t.Inputer.Start()
	go t.Outputer.Start()
	go func() {
		var err error
		for {
			if err = t.recv(); err != nil {
				t.logs.Error("Transport recv:%v", err)
				close(t.r_chan)
				break
				time.Sleep(1 * time.Second)
			}
		}
	}()
	go func() {
		for {
			if err := t.handle(); err != nil {
				t.logs.Error("Transport handle:%v", err)
				close(t.w_chan)
				break
			}
		}
	}()
	go func() {
		for {
			if err := t.send(); err != nil {
				t.logs.Error("Transport send:%v", err)
				t.Output.Close()
				t.isEnd <- true
				break
				time.Sleep(1 * time.Second)
			}
		}
	}()
	t.logs.Info("Transport start success...%s", time.Now())
	select {
	case <-t.isEnd:
		break
	}
}

func (t *Transport) Stop() {
	//t.Input.Close()
	t.Output.Close()
	t.logs.Info("Transport stop success...%s", time.Now())
}
