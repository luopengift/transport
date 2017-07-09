package file

import (
	"github.com/luopengift/golibs/file"
)

type Input struct {
	*file.Tail
}

func New(name string) *Input {
	input := new(Input)
	input.Tail = file.NewTail(name, file.NullRule)
	return input
}

func (in *Input) Read(p []byte) (int, error) {
	return 0, nil
}

func (in *Input) Close() error {
	return in.Tail.Close()
}
