package plugin_tcp

import (
	"code.google.com/p/gogoprotobuf/proto"
	"encoding/json"
	"fmt"
	"net"
	"stalin/message"
	. "stalin/plugins"
	"sync"
	"time"
)

func init() {
	RegisterPlugin("TcpInput", new(TcpInput))
}

type TcpInput struct {
	config   *TcpInputConfig
	gConfig  *GlobalConfig
	listener net.Listener
	wg       sync.WaitGroup
}

type TcpInputConfig struct {
	Address      string `json:"address"`
	Net          string `json:"net"`
	KeepAliveSec int    `json:"keep_alive"`
	keepAlive    time.Duration
	Name         string `json:"name"`
}

func (t *TcpInput) Init(g *GlobalConfig, name string) (Plugin, error) {
	plugin := &TcpInput{gConfig: g}
	// make default config
	config := &TcpInputConfig{
		Net:          "tcp",
		Name:         name,
		KeepAliveSec: 0,
		Address:      "0.0.0.0:5555",
	}
	if err := json.Unmarshal(g.PluginConfigs[name], config); err != nil {
		return nil, err
	}
	plugin.config = config
	return plugin, nil
}

func (t *TcpInput) Inject(msg *message.Message) error {
	return PluginNotImplemented
}

func (t *TcpInput) runListen() error {

	listener, err := net.Listen(t.config.Net, t.config.Address)
	if err != nil {
		return fmt.Errorf("[TcpInput]: ListenTCP failed: %s\n", err.Error())
	}

	t.listener = listener
	if t.config.Net == "tcp" && t.config.KeepAliveSec != 0 {
		t.config.keepAlive = time.Duration(t.config.KeepAliveSec) * time.Second
	}

	return nil
}

func (t *TcpInput) handleConnection(conn net.Conn) {

	defer func() {
		conn.Close()
		t.wg.Done()
	}()

	LogDebug("[TCPInput]: New client: %v", conn.RemoteAddr())

	for {
		msg := &message.Message{}
		data, err := message.Unpack(conn)
		if err != nil {
			LogErr("[TCPInput]: Decoding message from client: %v error: %v", conn.RemoteAddr(), err.Error())
			return
		}
		err = proto.Unmarshal(data, msg)
		if err != nil {
			LogErr("[TCPInput]: Unmarshal message from client: %v error: %v", conn.RemoteAddr(), err.Error())
			return
		}
		t.gConfig.Output(msg)
	}
}

func (t *TcpInput) Run() error {

	if err := t.runListen(); err != nil {
		return err
	}

	var conn net.Conn
	var e error

	LogInfo("[TCPInput]: Started at: %v", t.config.Address)

	for {

		if conn, e = t.listener.Accept(); e != nil {
			if e.(net.Error).Temporary() {
				LogErr("[TCPInput]: TCP accept failed: %s", e)
				continue
			} else {
				LogErr("[TCPInput]: TCP accept error: %s", e)
				break
			}
		}

		if t.config.KeepAliveSec != 0 {
			tcpConn, ok := conn.(*net.TCPConn)
			if !ok {
				return fmt.Errorf("KeepAlive only supported for TCP Connections.")
			}
			tcpConn.SetKeepAlive(true)
			tcpConn.SetKeepAlivePeriod(t.config.keepAlive)
		}

		t.wg.Add(1)
		go t.handleConnection(conn)

	}
	t.wg.Wait()

	return nil
}
