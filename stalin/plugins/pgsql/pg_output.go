package plugin_pgsql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"stalin/message"
	. "stalin/plugins"
)

func init() {
	RegisterPlugin("PgOutput", new(PgOutput))
}

type PgOutput struct {
	config  *PgOutputConfig
	gConfig *GlobalConfig
	db      *sql.DB
	Stat    *Counter
}

type PgOutputConfig struct {
	Name          string `json:"name"`
	Connection    string `json:"connection_string"`
	Query         string `"json:"query"`
	PoolSize      int    `json:"pool_size"`
	MaxQueue      int    `json:"max_queue"`
	Statistic     bool   `json:"statistic"`
	StatisticTime int    `json:"statistic_time"`
}

func (t *PgOutput) Init(g *GlobalConfig, name string) (Plugin, error) {
	plugin := &PgOutput{gConfig: g}
	// make default config
	config := &PgOutputConfig{
		Name:          name,
		Connection:    "postgres://riemann_face_user@127.0.0.1/riemann_face?sslmode=disable",
		Query:         "SELECT UpdateProblems($1, $2, $3, $4, $5, $6, $7)",
		PoolSize:      40,
		Statistic:     true,
		StatisticTime: 1,
	}
	if err := json.Unmarshal(g.PluginConfigs[name], config); err != nil {
		return nil, err
	}
	plugin.config = config
	plugin.Stat = NewCounter(name, config.StatisticTime, config.Statistic)
	plugin.Stat.Start()
	return plugin, nil
}

func (p *PgOutput) eventToPg(event *message.Event) {
	p.Stat.In()
	host, domain := event.GetHostDomain()
	// get metric
	var metric sql.NullFloat64
	if m, err := event.GetMetric(); err == nil {
		metric = sql.NullFloat64{m, true}
	}
	// get State
	State := &sql.NullString{String: event.GetState()}
	State.Valid = (State.String != "")
	rows, err := p.db.Query(p.config.Query,
		host,
		domain,
		event.GetService(),
		State,
		event.GetDescription(),
		metric,
		event.GetSQLTags(),
	)
	if err != nil {
		LogErr("[PGOut]: Query %v", err)
	} else {
		rows.Close()
	}
	go p.Stat.Out()
}

func (t *PgOutput) Inject(msg *message.Message) error {
	for _, event := range msg.GetEvents() {
		if t.Stat.QueueSize() > t.config.MaxQueue {
			t.Stat.UpDropped()
			return fmt.Errorf("Max queue size error")
		}
		t.eventToPg(event)
	}
	return nil
}

func (t *PgOutput) Run() error {
	db, err := sql.Open("postgres", t.config.Connection)
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(t.config.PoolSize)
	db.SetMaxIdleConns(t.config.PoolSize)
	db.Prepare(t.config.Query)
	t.db = db
	LogInfo("[PGOut]: Started with connection string: %v", t.config.Connection)
	return nil
}
