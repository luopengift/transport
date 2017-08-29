package transport

type PluginsMap struct {
	Inputers  map[string]Inputer
	Outputers map[string]Outputer
	Handlers  map[string]Handler
}

func NewPluginsMap() *PluginsMap {
	Plugins = new(PluginsMap)
	Plugins.Inputers = make(map[string]Inputer)
	Plugins.Outputers = make(map[string]Outputer)
	Plugins.Handlers = make(map[string]Handler)
	return Plugins
}

var Plugins *PluginsMap

func RegistInputer(key string, input Inputer) {
	Plugins.Inputers[key] = input
}

func RegistOutputer(key string, output Outputer) {
	Plugins.Outputers[key] = output
}

func RegistHandler(key string, handle Handler) {
	Plugins.Handlers[key] = handle
}

func init() {
	Plugins = NewPluginsMap()
}
