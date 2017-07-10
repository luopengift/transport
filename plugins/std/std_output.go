package plugins

import (
        "github.com/luopengift/transport"
        "os"
)

type Stdout struct {
        *os.File
}

func NewStdout() *Stdout {
        return new(Stdout)
}

func (stdout *Stdout) Close() error {
        return stdout.Close()
}

func (stdout *Stdout) Init(config map[string]string) error {
        stdout.File = os.Stdout
        return nil
}



func init() {
        transport.RegistOutputer("stdout", NewStdout())
}

