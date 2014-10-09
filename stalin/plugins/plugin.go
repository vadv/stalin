package plugins

import (
	"reflect"
	"stalin/message"
)

var AvailablePlugins = make(map[string]reflect.Type)

func RegisterPlugin(name string, typ interface{}) {
	AvailablePlugins[name] = reflect.TypeOf(typ)
}

type Plugin interface {
	Init(g *GlobalConfig, name string) (Plugin, error)
	Run() error
	Inject(msg *message.Message) error
}

func pluginFromName(name string) (Plugin, error) {
	typ := AvailablePlugins[name]
	if typ == nil {
		return nil, PluginNotRegister
	}
	v := reflect.New(typ).Elem()
	return v.Interface().(Plugin), nil
}
