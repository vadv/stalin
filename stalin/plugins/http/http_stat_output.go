package plugin_http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"stalin/message"
	. "stalin/plugins"
	"time"
)

func init() {
	RegisterPlugin("HttpStatOutput", new(HttpStatOutput))
}

type HttpStatOutput struct {
	gConfig       *GlobalConfig
	config        *HttpStatOutputConfig
	msg_total     int64
	event_total   int64
	msg_queue     int
	event_queue   int
	msg_per_sec   float64
	event_per_sec float64
}

type HttpStatOutputConfig struct {
	Address string `json:"address"`
	ApiUrl  string `json:"api_url"`
}

type apiOutput struct {
	MsgTotal       int64           `json:"msg_total"`
	EventTotal     int64           `json:"event_total"`
	MsgPerSecond   float64         `json:"msg_per_sec"`
	EventPerSecond float64         `json:"event_per_sec"`
	Plugins        []apiPluginInfo `json:"plugins"`
}

type apiPluginInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func (r *HttpStatOutput) Init(g *GlobalConfig, name string) (Plugin, error) {
	plugin := &HttpStatOutput{gConfig: g}
	// make default config
	config := &HttpStatOutputConfig{
		Address: "0.0.0.0:55855",
		ApiUrl:  "/status.json",
	}
	if err := json.Unmarshal(g.PluginConfigs[name], config); err != nil {
		return nil, err
	}
	plugin.config = config
	return plugin, nil
}

func (r *HttpStatOutput) Inject(msg *message.Message) error {
	r.msg_total++
	r.msg_queue++
	for _, _ = range msg.GetEvents() {
		r.event_total++
		r.event_queue++
	}
	return nil
}

func (r *HttpStatOutput) Run() error {
	LogInfo("[HttpStatOutput]: Started at http://%v%v", r.config.Address, r.config.ApiUrl)
	http.HandleFunc(r.config.ApiUrl, r.apiHandler)
	go r.tick()
	if err := http.ListenAndServe(r.config.Address, nil); err != nil {
		return err
	}
	return nil
}

func (r *HttpStatOutput) tick() {
	tickchan := time.Tick(time.Second * 5)
	for {
		select {
		case <-tickchan:
			r.msg_per_sec = float64(r.msg_queue) / 5
			r.event_per_sec = float64(r.event_queue) / 5
			r.event_queue = 0
			r.msg_queue = 0
		}
	}
}

func (p *HttpStatOutput) apiHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	pluginInfo := []apiPluginInfo{}
	for name, _ := range p.gConfig.Outputs {
		info := &apiPluginInfo{Name: name, Type: "Output"}
		pluginInfo = append(pluginInfo, *info)
	}
	for name, _ := range p.gConfig.Inputs {
		info := &apiPluginInfo{Name: name, Type: "Input"}
		pluginInfo = append(pluginInfo, *info)
	}

	out := &apiOutput{
		Plugins:        pluginInfo,
		MsgTotal:       p.msg_total,
		EventTotal:     p.event_total,
		MsgPerSecond:   p.msg_per_sec,
		EventPerSecond: p.event_per_sec,
	}
	out_byte, _ := json.Marshal(out)
	fmt.Fprintf(w, string(out_byte))

}
