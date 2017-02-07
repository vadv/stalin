package opentsdb

import (
	"net"

	"stalin/events"
	. "stalin/plugins"
)

func init() {
	Plugin.Creator(newOpentsdbOutput).Type("OpentsdbOutput").Description("Send external messages to opentsdb").Register()
}

type OpentsdbOutput struct {
	buffWriter *events.BuffWriter
	Log        *Logger `json:"-"`
	Address    string  `json:"address" description:"Address for connect"`
	Size       int     `json:"buffer" description:"Buffer size in bytes"`
	Flush      int     `json:"flush" description:"Time for periodic flush (in ms)"`
	Reconnect  int     `json:"reconnect" description:"Sleep time before reconnect (in ms)"`
	QueueSize  int     `json:"queue_size" description:"Max queue size"`
	errChan    chan error
}

func (t *OpentsdbOutput) Connect() (net.Conn, error) {
	return net.Dial("tcp", t.Address)
}

func (t *OpentsdbOutput) Inject(events *events.Events) error {
	for _, event := range events.List() {
		t.buffWriter.Inject(&buffWriterEvent{event})
	}
	return nil
}

func newOpentsdbOutput(name string) PluginInterface {
	return &OpentsdbOutput{
		Log:       NewLog(name),
		Address:   "127.0.0.1:4242",
		errChan:   make(chan error),
		Size:      1024 * 64,
		Flush:     500,
		Reconnect: 500,
		QueueSize: 64 * 1024,
	}
}

func (c *OpentsdbOutput) Start() error {

	c.buffWriter = events.NewBuffWriter(c.Connect, c.errChan, c.Flush, c.Reconnect, c.Size, c.QueueSize)
	go c.buffWriter.Run()

	c.Log.Info("Start for address: %v", c.Address)

	for {
		select {
		case err := <-c.errChan:
			c.Log.Errorf("Error: %v", err)
		}
	}

}
