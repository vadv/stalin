package plugins

import (
	"encoding/json"
	"github.com/bbangert/toml"
	"io/ioutil"
	"regexp"
	"stalin/message"
	"strings"
)

type GlobalConfig struct {
	Outputs       map[string]Plugin
	Inputs        map[string]Plugin
	Router        Plugin
	PluginConfigs map[string][]byte
	StalinConfig  *DaemonConfig
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
		StalinConfig:  RuntimeStalinConfig,
	}

	// пробегаемся по конфигу, заполняем Input и Ouput
	for name, conf := range configFile {

		configSection = nil
		if err := toml.PrimitiveDecode(conf, &configSection); err != nil {
			LogErr("Can't decode '%v', section: %#v", name, conf)
			return nil, err
		}

		if strings.ToLower(name) == "stalin" {
			if err := gConfig.StalinConfig.SetConfig(configSection); err != nil {
				LogErr("Can't set stalin config: %v", err)
			}
			continue
		}

		if err := toml.PrimitiveDecode(configSection["type"], &pluginType); err != nil {
			LogErr("Can't get type for section name: %v", name)
			continue
		}

		// пытаемся создать плагин
		plugin, err := pluginFromName(pluginType)
		if err != nil {
			LogErr("Init plugin '%v' fatal error: %v", pluginType, err)
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

var PluginTypeRegex = regexp.MustCompile("(Router|Input|Output)$")

func getPluginCategory(pluginType string) string {
	pluginCats := PluginTypeRegex.FindStringSubmatch(pluginType)
	if len(pluginCats) < 2 {
		return ""
	}
	return pluginCats[1]
}

func (g *GlobalConfig) Output(msg *message.Message) {
	go g.Router.Inject(msg)
}
