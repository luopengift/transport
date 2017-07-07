package transport

type Handler interface {
	Handle(in, out []byte) error
}

type Filter struct {
    Handler
}

func NewFilter(h Handler) *Filter {
    f := new(Filter)
    f.Handler = h
    return f
}

func (f *Filter) Handle(in,out []byte) error {
    return f.Handler.Handle(in,out)
}


