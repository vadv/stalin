package tcp

import (
	"net"

	"stalin/events"
	. "stalin/plugins"
)

func init() {
	Plugin.Creator(newTcpOutput).Type("TcpOutput").Description("Send external messages to address (use protobuf)").Register()
}

type TcpOutput struct {
	buffWriter *events.BuffWriter
	Log        *Logger `json:"-"`
	Address    string  `json:"address" description:"Address for connect"`
	Size       int     `json:"buffer" description:"Buffer size in bytes"`
	Flush      int     `json:"flush" description:"Time for periodic flush (in ms)"`
	Reconnect  int     `json:"reconnect" description:"Sleep time before reconnect (in ms)"`
	QueueSize  int     `json:"queue_size" description:"Max queue size"`
	errChan    chan error
}

func (t *TcpOutput) Connect() (net.Conn, error) {
	return net.Dial("tcp", t.Address)
}

func (t *TcpOutput) Inject(events *events.Events) error {
	t.buffWriter.Inject(&buffWriterEvent{Events: events})
	return nil
}

func newTcpOutput(name string) PluginInterface {
	return &TcpOutput{
		Log:       NewLog(name),
		Address:   "127.0.0.1:5555",
		Size:      1024 * 64,
		Flush:     500,
		Reconnect: 500,
		QueueSize: 64 * 1024,
		errChan:   make(chan error),
	}
}

func (t *TcpOutput) Start() error {

	t.buffWriter = events.NewBuffWriter(t.Connect, t.errChan, t.Flush, t.Reconnect, t.Size, t.QueueSize)
	go t.buffWriter.Run()

	t.Log.Info("Start for address: %v", t.Address)

	for {
		select {
		case err := <-t.errChan:
			t.Log.Errorf("Error: %v", err)
		}
	}
}
