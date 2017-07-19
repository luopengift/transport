package transport

import (
	"github.com/luopengift/golibs/logger"
	"time"
)

const (
	B = 1        //1B = 8bit
	K = 1024 * B //1KB
	M = 1024 * K //1MB
	G = 1024 * M //1GB

	MAX = 1 * M
)

type Transport struct {
	*Input
	*Output
	*Middleware
	ReadBuffer  chan []byte
	WriteBuffer chan []byte
}

func NewTransport(in Inputer, h Handler, out Outputer) *Transport {
	transport := new(Transport)
	transport.Input = NewInput(in)
	transport.Middleware = NewMiddleware(h)
	transport.Output = NewOutput(out)
	transport.ReadBuffer = make(chan []byte, 10)
	transport.WriteBuffer = make(chan []byte, 10)
	return transport

}

// 将数据从read接口读入 ReadBuffer中
func (t *Transport) recv() {
	t.Input.Mutex.Lock()
	defer t.Input.Mutex.Unlock()
	if t.Inputer == nil {
		logger.Debug("input is nil")
		time.Sleep(1000 * time.Millisecond)
		return
	}
	b := make([]byte, MAX, MAX)
	n, err := t.Inputer.Read(b)
	if err != nil {
		logger.Error("recv error:%v", err)
	}
	t.ReadBuffer <- b[:n]
	logger.Debug("recv %v", string(b[:n]))
}

func (t *Transport) handle() {
	tmp := <-t.ReadBuffer
	go func(p []byte) {
		b := make([]byte, MAX, MAX)
		n, err := t.Handler.Handle(p, b)
		if err != nil {
			logger.Error("Handler Error!%v", err)
			return
		}
		t.WriteBuffer <- b[:n]
	}(tmp)
}

// 将数据从WriteBuffer写入 Write接口中
func (t *Transport) send() {
	t.Output.Mutex.Lock()
	defer t.Output.Mutex.Unlock()
	if t.Outputer == nil {
		logger.Debug("output is nil")
		time.Sleep(1000 * time.Second)
		return
	}
	b := <-t.WriteBuffer
	logger.Debug("send %v", string(b))
	n, err := t.Outputer.Write(b)
	if err != nil {
		logger.Error("send error:%v,%v|message:", n, err, string(b))
	}
}

func (t *Transport) Run() {
	go t.Inputer.Start()
	go t.Outputer.Start()
	go func() {
		for {
			t.recv()
		}
	}()
	go func() {
		for {
			t.handle()
		}
	}()
	go func() {
		for {
			t.send()
		}
	}()
	logger.Info("Transport start success...%s\n", time.Now())
	select {}
}

func (t *Transport) Stop() {
	t.Input.Close()
	close(t.ReadBuffer)
	//t.Filter.Close()
	close(t.WriteBuffer)
	t.Output.Close()

}
