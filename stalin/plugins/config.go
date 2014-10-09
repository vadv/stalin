package plugins

import (
	"encoding/json"
	"github.com/bbangert/toml"
	"io/ioutil"
	"log"
	"regexp"
	"stalin/message"
)

type GlobalConfig struct {
	Outputs       map[string]Plugin
	Inputs        map[string]Plugin
	Router        Plugin
	PluginConfigs map[string][]byte
}

type PluginConfig map[string]toml.Primitive
type ConfigFile PluginConfig

func NewGlobalConfigFromFile(filename string) (*GlobalConfig, error) {

	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var configFile ConfigFile
	var configSection ConfigFile
	var pluginType string

	if _, err := toml.Decode(string(contents), &configFile); err != nil {
		return nil, err
	}

	gConfig := &GlobalConfig{
		Outputs:       make(map[string]Plugin),
		Inputs:        make(map[string]Plugin),
		PluginConfigs: make(map[string][]byte),
	}

	// пробегаемся по конфигу, заполняем Input и Ouput
	for name, conf := range configFile {

		configSection = nil
		if err := toml.PrimitiveDecode(conf, &configSection); err != nil {
			log.Printf("Can't decode '%v', section: %#v\n", name, conf)
			return nil, err
		}

		if err := toml.PrimitiveDecode(configSection["type"], &pluginType); err != nil {
			log.Printf("Can't get type for section name : %v\n", name)
			continue
		}

		// пытаемся создать плагин
		plugin, err := pluginFromName(pluginType)
		if err != nil {
			log.Printf("Init plugin '%v' fatal error: %v\n", pluginType, err)
			return nil, err
		} else {
			// плагин создался без ошибки
			// конвертим конфиг в json bytes
			bytes, err := json.Marshal(configSection)
			if err != nil {
				return nil, err
			}
			gConfig.PluginConfigs[name] = bytes
			// добавляем в gConfig по категориям
			switch getPluginCategory(pluginType) {
			case "Input":
				gConfig.Inputs[name] = plugin
			case "Output":
				gConfig.Outputs[name] = plugin
			}
		}
	}

	// перестраеваем gConfig: теперь нам известны все плагины, пытаемся их инициализировать
	for name, input := range gConfig.Inputs {
		plugin, err := input.Init(gConfig, name)
		if err != nil {
			return nil, err
		}
		gConfig.Inputs[name] = plugin
	}

	for name, input := range gConfig.Outputs {
		plugin, err := input.Init(gConfig, name)
		if err != nil {
			return nil, err
		}
		gConfig.Outputs[name] = plugin
	}

	// мы заполнили Input и Output плагины, добавляем Router
	route_interface, _ := pluginFromName("Router")
	router, _ := route_interface.Init(gConfig, "router")
	gConfig.Router = router

	return gConfig, nil
}

func (g *GlobalConfig) Output(msg *message.Message) {
	go g.Router.Inject(msg)
}

var PluginTypeRegex = regexp.MustCompile("(Router|Input|Output)$")

func getPluginCategory(pluginType string) string {
	pluginCats := PluginTypeRegex.FindStringSubmatch(pluginType)
	if len(pluginCats) < 2 {
		return ""
	}
	return pluginCats[1]
}
