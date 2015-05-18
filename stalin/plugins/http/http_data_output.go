package http

import (
	"bytes"
	"net/http"
	"time"

	"stalin/events"
	. "stalin/plugins"
)

func init() {
	Plugin.Creator(newHttpDataOutput).Type("HttpDataOutput").Description("Send data of external messages via http request").Register()
}

type HttpDataOutput struct {
	Log     *Logger      `json:"-"`
	Url     string       `json:"url" description:"Address of http receiver"`
	Timeout int          `json:"timeout" description:"Timeout for request (in sec)"`
	Queue   string       `json:"queue" description:"Data type of external messages"`
	Method  string       `json:"method" description:"Method of http-request"`
	Retry   int          `json:"retry" description:"Count of retry send data"`
	client  *http.Client `json:"-"`
}

func newHttpDataOutput(name string) PluginInterface {
	return &HttpDataOutput{
		Log:     NewLog(name),
		Url:     "http://localhost:9200",
		Method:  "POST",
		Timeout: 1,
		Retry:   1,
	}
}

func (p *HttpDataOutput) Inject(events *events.Events) error {
	if events.DataType != nil && events.Data != nil && *(events.DataType) == p.Queue {
		data := []byte(*(events.Data))
		req, err := http.NewRequest(p.Method, p.Url, bytes.NewBuffer(data))
		if err != nil {
			p.Log.Error("Make request: %#v", err.Error())
			return err
		}
		resp, err := p.client.Do(req)
		if err != nil {
			p.Log.Error("Send request: %#v", err.Error())
			return err
		}
		if resp.StatusCode > 300 {
			p.Log.Error("Status code: %d", resp.StatusCode)
		}
	}
	return nil
}

func (p *HttpDataOutput) Start() error {
	p.Log.Info("Start for address: %v", p.Url)
	p.client = &http.Client{
		Timeout: (time.Duration(p.Timeout) * time.Second),
	}
	return nil
}
