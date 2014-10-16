package plugin_carbon

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
	RegisterPlugin("CarbonOutput", new(CarbonOutput))
}

type CarbonOutput struct {
	config   *CarbonOutputConfig
	gConfig  *GlobalConfig
	conn     net.Conn
	ichan    chan []byte
	writer   *bufio.Writer
	tick     time.Duration
	tickchan <-chan time.Time
	Stat     *Counter
}

type CarbonOutputConfig struct {
	Name           string `json:"name"`
	Address        string `json:"address"`
	TickSec        int    `json:"flush_tick"`
	MaxQueue       int    `json:"max_queue"`
	BufferSize     int    `json:"buff_size"`
	MaxMessageSize int    `json:"max_message_size"`
	Statistic      bool   `json:"statistic"`
	StatisticTime  int    `json:"statistic_time"`
}

func (t *CarbonOutput) Init(g *GlobalConfig, name string) (Plugin, error) {
	plugin := &CarbonOutput{gConfig: g}
	// make default config
	config := &CarbonOutputConfig{
		Address:        "graphite-1.production.infra.home:2023",
		Name:           name,
		MaxQueue:       10000,
		BufferSize:     1024 * 64,
		MaxMessageSize: 1024 * 64,
		Statistic:      true,
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

func (t *CarbonOutput) connect() error {
	conn, err := net.Dial("tcp", t.config.Address)
	if err != nil {
		return err
	}
	t.conn = conn
	t.writer = bufio.NewWriterSize(t.conn, t.config.BufferSize)
	return nil
}

func (t *CarbonOutput) loopConnect() {
	for {
		if err := t.connect(); err != nil {
			LogErr("[CarbonOut]: Connect to %v error: %v", t.config.Address, err)
			time.Sleep(time.Second)
			continue
		}
		break
	}
}

func (t *CarbonOutput) reciveAndTick() {
	t.loopConnect()
	for {
		select {
		case <-t.tickchan:
			if err := t.writer.Flush(); err != nil {
				LogErr("[CarbonOut]: Flush: %v", err)
			}
		case data := <-t.ichan:
			if _, err := t.writer.Write(data); err != nil {
				t.loopConnect()
			}
		}
	}
}

func (t *CarbonOutput) Inject(msg *message.Message) error {
	for _, event := range msg.GetEvents() {
		if t.Stat.QueueSize() > t.config.MaxQueue {
			t.Stat.UpDropped()
			return fmt.Errorf("Max queue size error")
		}
		t.Stat.In()
		graphitemsg, err := event.ToGraphite()
		if err == nil {
			t.ichan <- []byte(graphitemsg)
		}
		t.Stat.Out()
	}
	return nil
}

func (t *CarbonOutput) Run() error {
	t.ichan = make(chan []byte, 10000)
	t.tickchan = time.Tick(t.tick)
	LogInfo("[CarbonOut]: Started with address %v", t.config.Address)
	t.reciveAndTick()
	return nil
}
