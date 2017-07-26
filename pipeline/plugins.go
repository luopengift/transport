package pipeline

var InputPlugins = map[string]Inputer{}

func RegistInputer(key string, in Inputer) {
	InputPlugins[key] = in
}

// 输出组件列表
var OutputPlugins = map[string]Outputer{}

func RegistOutputer(key string, out Outputer) {
	OutputPlugins[key] = out
}

var HandlePlugins = map[string]Handler{}

func RegistHandler(key string, h Handler) {
	HandlePlugins[key] = h
}
