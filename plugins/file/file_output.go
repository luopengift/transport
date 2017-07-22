package file

import (
	"github.com/luopengift/golibs/file"
	"github.com/luopengift/transport"
)

type FileOutput struct {
	*file.Tail
}

func NewFileOutput() *FileOutput {
	return new(FileOutput)
}

func (in *FileOutput) Init(cfg map[string]string) error {
	in.Tail = file.NewTail(cfg["path"], file.NullRule)
	return nil
}

func (in *FileOutput) Write(p []byte) (int, error) {
	return in.Tail.Read(p)
}

func (in *FileOutput) Start() error {
	in.Tail.ReadLine()
	return nil
}

func (in *FileOutput) Close() error {
	return nil //stdout.Close()
}

func init() {
	transport.RegistOutputer("file", NewFileOutput())
}
