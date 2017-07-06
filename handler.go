package transport

type Handler interface{
    Handle(in, out []byte) error
}
