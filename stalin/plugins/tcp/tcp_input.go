package tcp

import (
	"fmt"
	"net"
	"sync"
	"time"

	"stalin/events"
	. "stalin/plugins"
)

func init() {
	Plugin.Creator(newTcpInput).Type("TcpInput").Description("Receive protobuf messages").Register()
}

type TcpInput struct {
	Log          *Logger  `json:"-"`
	Clients      *Counter `json:"-"`
	Msg          *Counter `json:"-"`
	Address      string   `json:"address" description:"Listen address"`
	Net          string   `json:"net" description:"Net type"`
	KeepAliveSec int      `json:"keep_alive" description:"Keep alive for accepted connections (in sec)"`
	ls           net.Listener
	wg           *sync.WaitGroup
}

type TcpInputStatus struct {
	Clients  int64 `json:"client_count"`
	Messages int64 `json:"processed"`
}

func newTcpInput(name string) PluginInterface {
	return &TcpInput{
		Log:     NewLog(name),
		Clients: NewCounter(),
		Msg:     NewCounter(),
		Address: "0.0.0.0:5555", Net: "tcp", KeepAliveSec: 10,
	}
}

func (t *TcpInput) runListen() error {

	t.Log.Info("Start at address: %v", t.Address)

	ls, err := net.Listen(t.Net, t.Address)
	if err != nil {
		return fmt.Errorf("ListenTCP failed: %s", err.Error())
	}
	t.ls = ls
	return nil
}

func (t *TcpInput) setKeepAlive(conn net.Conn) error {
	if t.KeepAliveSec != 0 {
		tcpConn, ok := conn.(*net.TCPConn)
		if !ok {
			return fmt.Errorf("KeepAlive supported only for TCP Connections.")
		}
		tcpConn.SetKeepAlive(true)
		return tcpConn.SetKeepAlivePeriod(time.Duration(t.KeepAliveSec) * time.Second)
	}
	return nil
}

func (t *TcpInput) handleConn(conn net.Conn) {
	t.Log.Info("Handle new client connection for: %v", conn.RemoteAddr())

	t.Clients.Inc()
	defer func() {
		conn.Close()
		t.wg.Done()
		t.Clients.Dec()
	}()

	if err := t.setKeepAlive(conn); err != nil {
		t.Log.Error("KeepAlive set error: %v", err)
	}

	for {
		events, err := events.ReadConn(conn)
		if err != nil {
			t.Log.Error("While read error: %v", err)
			return
		}
		t.Msg.Inc()
		Pipeline.Inject(events)
	}

}

func (t *TcpInput) Start() error {

	t.wg = &sync.WaitGroup{}

	var conn net.Conn
	var e error

	if err := t.runListen(); err != nil {
		return err
	}
	for {
		if conn, e = t.ls.Accept(); e != nil {
			if e.(net.Error).Temporary() {
				continue
			} else {
				break
			}
		}
		t.wg.Add(1)
		go t.handleConn(conn)
	}

	t.wg.Wait()
	return nil
}

func (t *TcpInput) Status() (interface{}, error) {
	return &TcpInputStatus{
		Clients:  t.Clients.Count(),
		Messages: t.Msg.Count(),
	}, nil
}
