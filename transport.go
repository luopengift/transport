package transport

import (
	"github.com/luopengift/golibs/logger"
	"time"
)

var MaxBytes = 1000

func InitBytes(max int) []byte {
	return make([]byte, 0, max)
}

type Transport struct {
	*Input
	*Output
	*Filter
	ReadBuffer  chan *[]byte
	WriteBuffer chan *[]byte
}

func NewTransport(in Inputer, h Handler, out Outputer) *Transport {
	transport := new(Transport)
	transport.Input = NewInput(in)
	transport.Filter = NewFilter(h)
	transport.Output = NewOutput(out)
	transport.ReadBuffer = make(chan *[]byte, 10)
	transport.WriteBuffer = make(chan *[]byte, 10)
	return transport

}

// 将数据从read接口读入 ReadBuffer中
func (t *Transport) recv() {
	t.Input.Mutex.Lock()
	defer t.Input.Mutex.Unlock()
	if t.Inputer == nil {
		time.Sleep(600 * time.Millisecond)
		return
	}
	b := InitBytes(MaxBytes)
	_, err := t.Inputer.Read(b)
	logger.Debug("recv %v", string(b))
	if err != nil {
		logger.Error("recv error:%v", err)
	}
	t.ReadBuffer <- &b
}

// 将数据从WriteBuffer写入 Write接口中
func (t *Transport) send() {
	t.Output.Mutex.Lock()
	defer t.Output.Mutex.Unlock()
	if t.Outputer == nil {
		time.Sleep(1000 * time.Second)
		return
	}
	b := <-t.WriteBuffer
	logger.Debug("send %v", string(*b))
	n, err := t.Outputer.Write(*b)
	if err != nil {
		logger.Error("send error:%v,%v|message:", n, err, string(*b))
	}
}

func (t *Transport) handle() {
	b := InitBytes(MaxBytes)
	err := t.Handler.Handle(*(<-t.ReadBuffer), b)
	if err != nil {
		logger.Error("Handler Error!%v", err)
	}

	t.WriteBuffer <- &b
}

func (t *Transport) Run() {
	//	go t.Inputer.Start()
	//	go t.Outputer.Start()
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
	logger.Info("Transport start success...%s", time.Now())
	select {}
}

func (t *Transport) Stop() {
	t.Input.Close()
	close(t.ReadBuffer)
	//t.Filter.Close()
	close(t.WriteBuffer)
	t.Output.Close()

}
