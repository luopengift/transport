package transport

type Inject struct{}

func (i *Inject) InjectInput(p []byte) error {
	T.injectInput(p)
	return nil
}

func (i *Inject) InjectOutput(p []byte) error {
	T.injectOutput(p)
	return nil
}
