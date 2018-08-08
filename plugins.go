package transport

import (
	"fmt"
)

// PluginsMap Store available plugins, include input, adapt, output.
type PluginsMap struct {
	Inputers  map[string]Inputer
	Outputers map[string]Outputer
	Adapters  map[string]Adapter
}

// NewpluginsMap map instance
func NewPluginsMap() *PluginsMap {
	Plugins = new(PluginsMap)
	Plugins.Inputers = make(map[string]Inputer)
	Plugins.Outputers = make(map[string]Outputer)
	Plugins.Adapters = make(map[string]Adapter)
	return Plugins
}

// Plugins Global plugins instance
var Plugins *PluginsMap

// RegistInputer Regist input plugin
func RegistInputer(key string, input Inputer) {
	Plugins.Inputers[key] = input
}

// RegistOutputer Regist out plugin
func RegistOutputer(key string, output Outputer) {
	Plugins.Outputers[key] = output
}

// RegistHandler Regist adapt plugin
func RegistHandler(key string, a Adapter) {
	Plugins.Adapters[key] = a
}

// PluginDetail Show plugins information
func PluginDetail() string {
	str := fmt.Sprintf("%-16s %s\n", "[Inputs]", "version")
	for name, inputer := range Plugins.Inputers {
		str += fmt.Sprintf("  %-15s %s\n", name, inputer.Version())
	}
	str += fmt.Sprintf("%-17s\n", "[Adapters]")
	for name, adapter := range Plugins.Adapters {
		str += fmt.Sprintf("  %-15s %s\n", name, adapter.Version())
	}
	str += fmt.Sprintf("%-17s\n", "[Outputers]")
	for name, outputer := range Plugins.Outputers {
		str += fmt.Sprintf("  %-15s %s\n", name, outputer.Version())
	}
	return str
}

func init() {
	Plugins = NewPluginsMap()
}
