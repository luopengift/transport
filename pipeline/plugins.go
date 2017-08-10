package pipeline

type PluginsMap struct {
	Input  map[string]Inputer
	Output map[string]Outputer
	Handle map[string]Handler
}

func NewPluginsMap() *PluginsMap {
	pluginsMap = new(PluginsMap)
	pluginsMap.Input = make(map[string]Inputer)
	pluginsMap.Output = make(map[string]Outputer)
	pluginsMap.Handle = make(map[string]Handler)
	return pluginsMap
}

var pluginsMap *PluginsMap

func RegistInputer(key string, input Inputer) {
	pluginsMap.Input[key] = input
}

func RegistOutputer(key string, output Outputer) {
	pluginsMap.Output[key] = output
}

func RegistHandler(key string, handle Handler) {
	pluginsMap.Handle[key] = handle
}

func init() {
	pluginsMap = NewPluginsMap()
}
