package transport

type Handler interface {
	Handle(in, out []byte) (n int, err error)
}

type Middleware struct {
	Handler
}

func NewMiddleware(h Handler) *Middleware {
	f := new(Middleware)
	f.Handler = h
	return f
}

func (f *Middleware) Handle(in, out []byte) (int, error) {
	return f.Handler.Handle(in, out)
}

var HandlePlugins = map[string]Handler{}

func RegistHandler(key string, h Handler) {
	HandlePlugins[key] = h
}
