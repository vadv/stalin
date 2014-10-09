package plugin_opentsdb

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"stalin/message"
	. "stalin/plugins"
	"time"
)

func init() {
	RegisterPlugin("OpentsdbOutput", new(OpentsdbOutput))
}

type OpentsdbOutput struct {
	config   *OpentsdbOutputConfig
	gConfig  *GlobalConfig
	conn     net.Conn
	ichan    chan []byte
	writer   *bufio.Writer
	tick     time.Duration
	tickchan <-chan time.Time
	Stat     *Counter
}

type OpentsdbOutputConfig struct {
	Name           string `json:"name"`
	Address        string `json:"address"`
	TickSec        int    `json:"tick"`
	MaxQueue       int    `json:"max_queue"`
	BufferSize     int    `json:"buff_size"`
	MaxMessageSize int    `json:"max_message_size"`
	Statistic      bool   `json:"statistic"`
	StatisticTime  int    `json:"statistic_time"`
}

func (t *OpentsdbOutput) Init(g *GlobalConfig, name string) (Plugin, error) {
	plugin := &OpentsdbOutput{gConfig: g}
	// make default config
	config := &OpentsdbOutputConfig{
		Address:        "opentsdb.production.infra.home:4242",
		Name:           name,
		MaxQueue:       10000,
		BufferSize:     1024 * 64,
		MaxMessageSize: 1024 * 64,
		StatisticTime:  1,
		TickSec:        1,
	}
	if err := json.Unmarshal(g.PluginConfigs[name], config); err != nil {
		return nil, err
	}
	plugin.tick = time.Duration(config.TickSec) * time.Second
	plugin.config = config
	plugin.Stat = NewCounter(name, config.StatisticTime, config.Statistic)
	plugin.Stat.Start()
	return plugin, nil
}

func (t *OpentsdbOutput) connect() error {
	conn, err := net.Dial("tcp", t.config.Address)
	if err != nil {
		return err
	}
	t.conn = conn
	t.writer = bufio.NewWriterSize(t.conn, t.config.BufferSize)
	return nil
}

func (t *OpentsdbOutput) loopConnect() {
	for {
		if err := t.connect(); err != nil {
			fmt.Printf("[OpentsdbOut]: Error connect to %v error: %v\n", t.config.Address, err)
			time.Sleep(time.Second)
			continue
		}
		break
	}
}

func (t *OpentsdbOutput) reciveAndTick() {
	t.loopConnect()
	for {
		select {
		case <-t.tickchan:
			if err := t.writer.Flush(); err != nil {
				fmt.Printf("[OpentsdbOut]: Error to flush: %v\n", err)
			}
		case data := <-t.ichan:
			if _, err := t.writer.Write(data); err != nil {
				t.loopConnect()
			}
		}
	}
}

func (t *OpentsdbOutput) Inject(msg *message.Message) error {
	for _, event := range msg.GetEvents() {
		if t.Stat.QueueSize() > t.config.MaxQueue {
			t.Stat.UpDropped()
			return fmt.Errorf("Max queue size error")
		}
		t.Stat.In()
		tsdbmsg, err := event.ToTsdb()
		if err == nil {
			t.ichan <- []byte(tsdbmsg)
		}
		t.Stat.Out()
	}
	return nil
}

func (t *OpentsdbOutput) Run() error {
	t.ichan = make(chan []byte, 10000)
	t.tickchan = time.Tick(t.tick)
	fmt.Printf("[OpentsdbOut]: name: %v with addr:%v\n", t.config.Name, t.config.Address)
	t.reciveAndTick()
	return nil
}
