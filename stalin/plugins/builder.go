package plugins

import (
	"fmt"
	"github.com/lann/builder"
	"sort"

	"stalin/events"
)

type PluginInterface interface {
	Start() error
}

type PluginOutputInterface interface {
	PluginInterface
	Inject(*events.Events) error
}

type pluginCreator struct {
	Type        string
	Description string
	Creator     func(string) PluginInterface
}

type plugins map[string]pluginCreator

var availablePlugins = make(plugins)
var Plugin = builder.Register(pluginBuilder{}, pluginCreator{}).(pluginBuilder)

type pluginBuilder builder.Builder

func (b pluginBuilder) Type(elm string) pluginBuilder {
	return builder.Set(b, "Type", elm).(pluginBuilder)
}

func (b pluginBuilder) Description(elm string) pluginBuilder {
	return builder.Set(b, "Description", elm).(pluginBuilder)
}

func (b pluginBuilder) Creator(elm func(string) PluginInterface) pluginBuilder {
	return builder.Set(b, "Creator", elm).(pluginBuilder)
}

func (b pluginBuilder) Register() {
	creator := builder.GetStruct(b).(pluginCreator)
	availablePlugins[creator.Type] = creator
}

func GetPluginFromType(pluginType, name string) (PluginInterface, error) {
	plugin, ok := availablePlugins[pluginType]
	if !ok {
		return nil, fmt.Errorf("Plugin type %s not registered", pluginType)
	}
	if plugin.Creator == nil {
		return nil, fmt.Errorf("Plugin type %s has no 'Creator' field", pluginType)
	}
	return plugin.Creator(name), nil
}

func (p plugins) getSortedTypes() []string {
	slice := sort.StringSlice{}
	for name, _ := range p {
		slice = append(slice, name)
	}
	slice.Sort()
	return slice
}
