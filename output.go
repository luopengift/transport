package transport

import (
	"io"
	"sync"
)

// 数据输入接口, 实现了标准io库中的ReadCloser接口
type Outputer interface {
	// Writer 接口包装了基本的 Write 方法，用于将数据存入自身。
	// Write 方法用于将 p 中的数据写入到对象的数据流中，
	// 返回写入的字节数和遇到的错误。
	// 如果 p 中的数据全部被写入，则 err 应该返回 nil。
	// 如果 p 中的数据无法被全部写入，则 err 应该返回相应的错误信息。
	io.WriteCloser //Write(p []byte) (n int, err error), Close() error
	Start() error
}

type Output struct {
	Outputer
	*sync.Mutex
	//	IsSend bool
}

func NewOutput(out Outputer) *Output {
	o := new(Output)
	o.Outputer = out
	o.Mutex = new(sync.Mutex)
	//	o.IsSend = false
	return o
}

/*
func (o *Output) StopWrite() {
	o.IsSend = false
}

func (o *Output) StartWrite() {
	o.IsSend = true
}
*/
func (o *Output) Set(out Outputer) error {
	o.Outputer = out
	return nil
}

func (o *Output) Write(p []byte) (int, error) {
	return o.Outputer.Write(p)
}

func (o *Output) Start() error {
	return o.Outputer.Start()
}
func (o *Output) Close() error {
	return o.Outputer.Close()
}


// 输出组件列表
var OutputPlugins = map[string]Outputer{}

func RegistOutputer(key string, out Outputer) {
    OutputPlugins[key] = out
}


