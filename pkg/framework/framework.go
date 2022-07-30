package framework

type Plugin interface {
	Start()
}

var plugins = make(map[string]Plugin)

func Start() {
	for _, plugin := range plugins {
		plugin.Start()
	}
}

func RegisterPlugin(name string, plugin Plugin) {
	plugins[name] = plugin
}
