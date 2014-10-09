package plugin_http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"stalin/message"
	. "stalin/plugins"
)

func init() {
	RegisterPlugin("HttpPlayerJsonInput", new(HttpPlayerJsonInput))
}

type HttpPlayerJsonInput struct {
	gConfig *GlobalConfig
	config  *HttpPlayerJsonInputConfig
}

type HttpPlayerJsonInputConfig struct {
	Address     string `json:"address"`
	ApiUrl      string `json:"api_url"`
	ExtDataType string `json:"queue"`
}

type PlayerPayload struct {
	ErrorType   string `json:"type"`
	Timestamp   int    `json:"timestamp"`
	Endpoint    string `json:"endpoint"`
	Domain      string `json:"domain"`
	Description string `json:"error"`
	ErrorCode   int    `json:"code"`
	EntryPoint  string `json:"entry_point"`
}

func (r *HttpPlayerJsonInput) Init(g *GlobalConfig, name string) (Plugin, error) {
	plugin := &HttpPlayerJsonInput{gConfig: g}
	config := &HttpPlayerJsonInputConfig{
		Address:     "0.0.0.0:9090",
		ApiUrl:      "/",
		ExtDataType: "HttpPlayerJsonInput",
	}
	if err := json.Unmarshal(g.PluginConfigs[name], config); err != nil {
		return nil, err
	}
	plugin.config = config
	return plugin, nil
}

func (r *HttpPlayerJsonInput) Inject(msg *message.Message) error {
	return PluginNotImplemented
}

func (r *HttpPlayerJsonInput) Run() error {
	fmt.Printf("[HttpPlayerJsonInput]: start at: http://%v%v\n", r.config.Address, r.config.ApiUrl)
	http.HandleFunc(r.config.ApiUrl, r.apiHandler)
	if err := http.ListenAndServe(r.config.Address, nil); err != nil {
		return err
	}
	return nil
}

func (h *HttpPlayerJsonInput) apiHandler(wr http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		wr.WriteHeader(http.StatusNotAcceptable)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		wr.WriteHeader(http.StatusInternalServerError)
		return
	}
	playerPayload := &PlayerPayload{}
	if err := json.Unmarshal(body, playerPayload); err != nil {
		wr.WriteHeader(http.StatusNotAcceptable)
		return
	}
	// проводим простую валидацию
	if playerPayload.ErrorType == "" || playerPayload.ErrorCode == 0 {
		wr.WriteHeader(http.StatusNotAcceptable)
		return
	}
	// очищаем от лишних полей
	data, _ := json.Marshal(playerPayload)
	msg := message.NewDataMessage(h.config.ExtDataType, string(data))
	// отсылаем в роутер сообщений
	h.gConfig.Output(msg)
	wr.WriteHeader(http.StatusCreated)
}
