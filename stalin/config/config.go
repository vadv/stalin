package config

import (
	"encoding/json"
	"fmt"
	"github.com/bbangert/toml"
	"io/ioutil"
	. "stalin/plugins"
	"strings"
)

type PluginConfig map[string]toml.Primitive
type ConfigFile PluginConfig

func DescriptionForPlugin(name string) string {
	return PluginDescription(name)
}

func ReadConfig(file string) (*PluginPipeline, error) {

	contents, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var configFile ConfigFile
	var configSection ConfigFile
	var pluginType string

	if _, err := toml.Decode(string(contents), &configFile); err != nil {
		return nil, err
	}

	for name, conf := range configFile {

		configSection = nil
		if err := toml.PrimitiveDecode(conf, &configSection); err != nil {
			return nil, fmt.Errorf("Can't decode '%v', section: %#v, error: %v", name, conf, err)
		}

		if strings.ToLower(name) == "stalin" {
			if err := Stalin.Configure(configSection); err != nil {
				return nil, fmt.Errorf("Can't set stalin config: %v", err)
			}
			continue
		}

		if err := toml.PrimitiveDecode(configSection["type"], &pluginType); err != nil {
			return nil, fmt.Errorf("Can't get type for section name: %v, error: %v", name, err)
		}

		plugin, err := GetPluginFromType(pluginType, name)
		if err != nil {
			return nil, err
		}
		data, err := json.Marshal(configSection)
		if err != nil {
			return nil, fmt.Errorf("Can't unmarshal to json: %v for plugin type: %s", err, pluginType)
		}
		if err := json.Unmarshal(data, plugin); err != nil {
			return nil, fmt.Errorf("Can't set config for plugin: %s", name)
		}

		if err := Pipeline.AddPlugin(name, plugin); err != nil {
			return nil, err
		}
	}

	return Pipeline, nil

}
