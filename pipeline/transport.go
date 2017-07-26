package pipeline

import (
	"fmt"
	"github.com/luopengift/golibs/logger"
	"time"
)

type Transport struct {
	*Input
	*Output
	*Middleware
	ReadBuffer  chan []byte
	WriteBuffer chan []byte
	isEnd       chan bool
}

func NewTransport(in Inputer, h Handler, out Outputer) *Transport {
	transport := new(Transport)
	transport.Input = NewInput(in)
	transport.Middleware = NewMiddleware(h, 2000)
	transport.Output = NewOutput(out)
	transport.ReadBuffer = make(chan []byte, 1000)
	transport.WriteBuffer = make(chan []byte, 1000)
	transport.isEnd = make(chan bool)
	return transport

}

func (t *Transport) End() chan bool {
	return t.isEnd
}

// 将数据从read接口读入 ReadBuffer中
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
	t.ReadBuffer <- b[:n]
	logger.Debug("recv %v", string(b[:n]))
	return nil
}

func (t *Transport) handle() error {
	tmp, ok := <-t.ReadBuffer
	if !ok {
		return ReadBufferClosedError
	}
	t.Middleware.Channel.Add()
	go func(p []byte) {
		defer t.Middleware.Channel.Done()
		b := make([]byte, MAX, MAX)
		n, err := t.Handler.Handle(p, b)
		if err != nil {
			logger.Error("Handler Error!%v", err)
			return
		}
		t.WriteBuffer <- b[:n]
	}(tmp)
    return nil
}

// 将数据从WriteBuffer写入 Write接口中
func (t *Transport) send() error {
	t.Output.Mutex.Lock()
	defer t.Output.Mutex.Unlock()
	if t.Outputer == nil {
		return OutputNullError
	}
	write, ok := <-t.WriteBuffer
	if !ok {
		return fmt.Errorf("WriteBuffer is closed")
	}
    go func(b []byte) error {
	    n, err := t.Outputer.Write(b)
	    if err != nil {
		    return fmt.Errorf("send error:%v,%v|message:%s", n, err, string(b))
	    }
	    logger.Debug("send:%v|msg:%v", n, string(b))
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
                logger.Error("Transport recv:%v", err)
				close(t.ReadBuffer)
                break
				time.Sleep(1 * time.Second)
			}
		}
	}()
	go func() {
		for {
		    if err := t.handle(); err != nil {
                logger.Error("Transport handle:%v", err)
				close(t.WriteBuffer)
                break
            }
		}
	}()
	go func() {
		for {
			if err := t.send(); err != nil {
				logger.Error("Transport send:%v", err)
				t.Output.Close()
                t.isEnd <- true
                break
                time.Sleep(1 * time.Second)
			}
		}
	}()
	logger.Info("Transport start success...%s", time.Now())
    select {
        case <- t.isEnd:
            break
    }
}

func (t *Transport) Stop() {
	//t.Input.Close()
	t.Output.Close()
	logger.Info("Transport stop success...%s", time.Now())
}
