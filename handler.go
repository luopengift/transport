package transport

type Handler interface {
	Handle(in, out []byte) (n int,err error)
}

type Filter struct {
	Handler
}

func NewFilter(h Handler) *Filter {
	f := new(Filter)
	f.Handler = h
	return f
}

func (f *Filter) Handle(in, out []byte) (int, error) {
	return f.Handler.Handle(in, out)
}

var FilterPlugins = map[string]Handler{}

func RegistHandler(key string, h Handler) {
	FilterPlugins[key] = h
}
