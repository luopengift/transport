package transport

import "io"

// 数据输入接口, 实现了标准io库中的ReadCloser接口
type Outputer interface {
	// Writer 接口包装了基本的 Write 方法，用于将数据存入自身。
	// Write 方法用于将 p 中的数据写入到对象的数据流中，
	// 返回写入的字节数和遇到的错误。
	// 如果 p 中的数据全部被写入，则 err 应该返回 nil。
	// 如果 p 中的数据无法被全部写入，则 err 应该返回相应的错误信息。
	io.WriteCloser //Write(p []byte) (n int, err error), Close() error
}

type Output struct{}

func (o *Output) Write(p []byte) (int, error) {
	return len(p), nil
}

func (o *Output) Close() error {
	return nil
}
