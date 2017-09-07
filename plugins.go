package transport

type PluginsMap struct {
	Inputers  map[string]Inputer
	Outputers map[string]Outputer
	Adapters  map[string]Adapter
}

func NewPluginsMap() *PluginsMap {
	Plugins = new(PluginsMap)
	Plugins.Inputers = make(map[string]Inputer)
	Plugins.Outputers = make(map[string]Outputer)
	Plugins.Adapters = make(map[string]Adapter)
	return Plugins
}

var Plugins *PluginsMap

func RegistInputer(key string, input Inputer) {
	Plugins.Inputers[key] = input
}

func RegistOutputer(key string, output Outputer) {
	Plugins.Outputers[key] = output
}

func RegistHandler(key string, a Adapter) {
	Plugins.Adapters[key] = a
}

func init() {
	Plugins = NewPluginsMap()
}
