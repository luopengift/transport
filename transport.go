package transport

import (
    "time"
	"github.com/luopengift/golibs/logger"
)

var MaxBytes = 1000


type Transport struct {
	Inputer
    Handler
	Outputer
	ReadBuffer  chan *[]byte
	WriteBuffer chan *[]byte
    IsSend bool
}

func NewTransport(in Inputer, h Handler, out Outputer) *Transport {
	transport := new(Transport)
	transport.Inputer = in
    transport.Handler = h
	transport.Outputer = out
	transport.ReadBuffer = make(chan *[]byte, 10)
	transport.WriteBuffer = make(chan *[]byte, 10)
    transport.IsSend = false
    return transport

}

func (t *Transport) StopWrite() {
    t.IsSend = false
}

func (t *Transport) StartWrite() {
    t.IsSend = true
}

func (t *Transport) SetInputer(in Inputer) {
    t.Inputer = in
}

func (t *Transport) SetOutputer(out Outputer) {
    t.Outputer = out
}

// 将数据从read接口读入 ReadBuffer中
func (t *Transport) recv() {
	b := make([]byte, MaxBytes)
	_, err := t.Inputer.Read(b)
	if err != nil {
		logger.Error("recv error:%v", err)
	}
	t.ReadBuffer <- &b
}

// 将数据从WriteBuffer写入 Write接口中
func (t *Transport) send() {
	if !t.IsSend {
        logger.Warn("stop send,please wait...")
        time.Sleep(100* time.Millisecond)
        return
    }
    b := <-t.WriteBuffer
	logger.Debug("send %v,output %#v", string(*b),t.Outputer)
    n, err := t.Outputer.Write(*b)
	if err != nil {
		logger.Error("send error:%v,%v|message:", n, err,string(*b))
	}
}

func (t *Transport) handle() {
    b := make([]byte,MaxBytes)
    err := t.Handler.Handle(*(<-t.ReadBuffer),b)
    if err != nil {
        logger.Error("Handler Error!%v",err)
    }
    
    t.WriteBuffer <-&b
}

func (t *Transport) Run() {
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
    logger.Info("Transport start success...")
	select {}
}
