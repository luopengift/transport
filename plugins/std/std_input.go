package plugins

import (
    "os"
    "github.com/luopengift/transport"
)

var (
	Stdin = os.Stdin
)


func init() {
    transport.RegistInputer("stdin",Stdin)
}
