package tcp

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"time"

	"stalin/events/js"
	. "stalin/plugins"
)

func init() {
	Plugin.Creator(newNDJsonInput).Type("NDJsonInput").Description("Newline delimited JSON input").Register()
}

type NDJsonInput struct {
	Log          *Logger `json:"-"`
	Address      string  `json:"address" description:"Listen address"`
	Net          string  `json:"net" description:"Net type"`
	KeepAliveSec int     `json:"keep_alive" description:"Keep alive for accepted connections (in sec)"`
	ls           net.Listener
	wg           *sync.WaitGroup
}

func newNDJsonInput(name string) PluginInterface {
	return &NDJsonInput{
		Log:     NewLog(name),
		Address: "127.0.0.1:1555", Net: "tcp", KeepAliveSec: 10,
	}
}

func (t *NDJsonInput) runListen() error {

	t.Log.Info("Start at address: %v", t.Address)

	ls, err := net.Listen(t.Net, t.Address)
	if err != nil {
		return fmt.Errorf("ListenTCP failed: %s", err.Error())
	}
	t.ls = ls
	return nil
}

func (t *NDJsonInput) setKeepAlive(conn net.Conn) error {
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

func (t *NDJsonInput) handleConn(conn net.Conn) {
	t.Log.Info("Handle new client connection for: %v", conn.RemoteAddr())

	reader := bufio.NewReader(conn)
	scanner := bufio.NewScanner(reader)

	defer func() {
		conn.Close()
		t.wg.Done()
	}()

	if err := t.setKeepAlive(conn); err != nil {
		t.Log.Error("KeepAlive set: %v", err)
	}

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			t.Log.Error("While read: %v", err)
			return
		} else {
			events, err := js.EventsFromJson(scanner.Bytes())
			if err != nil {
				t.Log.Error("While unmarshal: %v", err)
				return
			}
			Pipeline.Inject(events)
		}
	}

}

func (t *NDJsonInput) Start() error {

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
