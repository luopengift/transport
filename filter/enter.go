package filter

import "bytes"

var enter = []byte("\n")
var none = []byte("")
var AddEnter *AddEnterSymbol = new(AddEnterSymbol)


// add a enter symbol at end of line, classic written into file
type AddEnterSymbol struct{}
func (a *AddEnterSymbol) Handle(in, out []byte) error {
    tmp := bytes.Join([][]byte{in,enter},[]byte{})
    copy(out,tmp)
    return nil
}


// direct connect input and output, do nothing
type DefaultConnection struct {}

func (d *DefaultConnection) Handle(in, out []byte) error {
    copy(out, in)
    return nil
}
