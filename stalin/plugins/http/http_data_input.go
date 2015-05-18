package http

import (
	"io/ioutil"
	"net/http"

	"stalin/events"
	. "stalin/plugins"
)

func init() {
	Plugin.Creator(newHttpDataInput).Type("HttpDataInput").Description("Receive http body request and send it's to external message").Register()
}

type HttpDataInput struct {
	Log   *Logger `json:"-"`
	Bind  string  `json:"bind" description:"Listen address"`
	Queue string  `json:"queue" description:"Data type of external messages"`
}

func newHttpDataInput(name string) PluginInterface {
	return &HttpDataInput{
		Log:  NewLog(name),
		Bind: "0.0.0.0:9090",
	}
}

func (p *HttpDataInput) Start() error {
	p.Log.Info("Start at address: %v", p.Bind)
	http.HandleFunc("/", p.apiHandler)
	if err := http.ListenAndServe(p.Bind, nil); err != nil {
		return err
	}
	return nil
}

func (p *HttpDataInput) apiHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	msg := events.NewEvents()
	data := string(body)
	msg.Data = &data
	msg.DataType = &(p.Queue)
	Pipeline.Inject(msg)
	w.WriteHeader(http.StatusOK)
	return
}
