package transport

import "io"

// 数据输入接口, 实现了标准io库中的ReadCloser接口
type Inputer interface {
	// Reader 接口包装了基本的 Read 方法，用于输出自身的数据。
	// Read 方法用于将对象的数据流读入到 p 中，返回读取的字节数和遇到的错误。
	// 在没有遇到读取错误的情况下：
	// 1、如果读到了数据（n > 0），则 err 应该返回 nil。
	// 2、如果数据被读空，没有数据可读（n == 0），则 err 应该返回 EOF。
	// 如果遇到读取错误，则 err 应该返回相应的错误信息。
	io.ReadCloser //Read(p []byte) (n int, err error),  Close() error
	//Read() ([]byte, error)
}

type Input struct{}

func (i *Input) Read(p []byte) (int, error) {
	return len(p), nil
}

func (i *Input) Close() error {
	return nil
}
