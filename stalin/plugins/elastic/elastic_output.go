package plugin_elastic

import (
	"encoding/json"
	"fmt"
	"github.com/mattbaird/elastigo/api"
	"github.com/mattbaird/elastigo/core"
	"stalin/message"
	. "stalin/plugins"
)

func init() {
	RegisterPlugin("ElasticOutput", new(ElasticOutput))
}

type ElasticOutput struct {
	gConfig *GlobalConfig
	config  *ElasticOutputConfig
}

type ElasticOutputConfig struct {
	Host         string `json:"host"`
	Verbose      bool   `json:"verbose"`
	ElasticIndex string `json:"elastic_index"`
	ElasticType  string `json:"elastic_type"`
	ExtDataType  string `json:"queue"`
}

func (r *ElasticOutput) Init(g *GlobalConfig, name string) (Plugin, error) {
	plugin := &ElasticOutput{gConfig: g}
	config := &ElasticOutputConfig{
		Host:         "localhost",
		ElasticIndex: "playermsg",
		ElasticType:  "msg",
		ExtDataType:  "HttpPlayerJsonInput",
	}
	if err := json.Unmarshal(g.PluginConfigs[name], config); err != nil {
		return nil, err
	}
	plugin.config = config
	return plugin, nil
}

func (r *ElasticOutput) Inject(msg *message.Message) error {
	if msg.GetDataType() == r.config.ExtDataType {
		if _, err := core.Index(r.config.ElasticIndex, r.config.ElasticType, "", nil, []byte(msg.GetData())); err != nil {
			return err
		}
	}
	return nil
}

func (r *ElasticOutput) Run() error {
	api.Domain = r.config.Host
	core.VerboseLogging = r.config.Verbose
	fmt.Printf("[ElasticOutput]: start with %v:%v\n", r.config.Host, r.config.Port)
	return nil
}
