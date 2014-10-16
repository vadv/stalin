package plugin_store

import (
	"encoding/json"
	"stalin/message"
	. "stalin/plugins"
	"stalin/plugins/store/api"
	"stalin/plugins/store/memory"
	. "stalin/plugins/store/problem"
	"time"
)

func init() {
	RegisterPlugin("MemStorePluginOutput", new(MemStorePluginOutput))
}

type MemStorePluginOutput struct {
	gConfig    *GlobalConfig
	config     *MemStorePluginOutputConfig
	store      *memory.Store
	savePeriod time.Duration
	defaultTtl int64 // sec
}

type MemStorePluginOutputConfig struct {
	FileName         string `json:"filename"`
	PeriodSave       int    `json:"save_period"`
	DefaultTtl       int64  `json:"default_ttl"`
	JanitorTtlPeriod int    `json:"janitor_ttl"`
	ApiAddress       string `json:"address"`
}

func (m *MemStorePluginOutput) Init(g *GlobalConfig, name string) (Plugin, error) {

	plugin := &MemStorePluginOutput{gConfig: g}
	config := &MemStorePluginOutputConfig{
		FileName:   "/tmp/store.gob",
		PeriodSave: 30, // 30 sec
		ApiAddress: "127.0.0.1:55455",
		DefaultTtl: 15 * 60, // 15min
	}
	if err := json.Unmarshal(g.PluginConfigs[name], config); err != nil {
		return nil, err
	}
	store := memory.NewStore(time.Second*time.Duration(config.DefaultTtl),
		time.Second*time.Duration(config.JanitorTtlPeriod))
	store.FileName = config.FileName

	now := time.Now()
	LogInfo("[MemStore]: Loading db: %v", store.FileName)
	if err := store.Init(); err != nil {
		return nil, err
	}
	LogInfo("[MemStore]: Loaded db: %v, spend time: %v", store.FileName, time.Now().Sub(now))

	plugin.config = config
	plugin.defaultTtl = config.DefaultTtl
	plugin.store = store
	plugin.savePeriod = time.Second * time.Duration(config.PeriodSave)
	// api
	httpApi := &api.Api{Address: config.ApiAddress, Store: store}
	go httpApi.StartApi()
	//
	return plugin, nil
}

func (m *MemStorePluginOutput) Inject(msg *message.Message) error {
	for _, event := range msg.Events {
		problem, err := NewProblemFromEvent(event, m.config.DefaultTtl)
		if err != nil {
			continue
		}
		m.store.SetProblem(problem)
	}
	return nil
}

func (m *MemStorePluginOutput) Run() error {
	LogInfo("[MemStore]: Started at: %v", m.config.ApiAddress)
	go m.tickAndSave()
	return nil
}

func (m *MemStorePluginOutput) tickAndSave() {
	tickchan := time.Tick(m.savePeriod)
	for {
		select {
		case <-tickchan:
			m.store.Save()
		}
	}
}
