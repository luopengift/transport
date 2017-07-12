package filter


var enter = []byte("\n")
var none = []byte("")
var AddEnter *AddEnterSymbol = new(AddEnterSymbol)

// add a enter symbol at end of line, classic written into file
type AddEnterSymbol struct{}

func (a *AddEnterSymbol) Handle(in, out []byte) (int,error) {
	in = append(in,'\n')
	n := copy(out, in)
    
    return n, nil
}

// direct connect input and output, do nothing
type DefaultConnection struct{}

func (d *DefaultConnection) Handle(in, out []byte) (int,error) {
    n := copy(out, in)
    return n, nil
}



