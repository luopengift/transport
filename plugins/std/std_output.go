package plugins

import (
    "os"
    "github.com/luopengift/transport"
)

var (
	Stdout = os.Stdout
)


func init() {
    transport.RegistOutputer("stdout",Stdout)
}
