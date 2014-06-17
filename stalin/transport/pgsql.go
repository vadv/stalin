package transport

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"stalin/riemann"
)

type PgTransport struct {
	e     chan []byte
	query string
	db    *sql.DB
	stat  *Stat
}

func NewPG(constring, query string, pool int, s *Stat) (*PgTransport, error) {

	db, err := sql.Open("postgres", constring)
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(pool)
	db.SetMaxIdleConns(pool)
	db.Prepare(query)

	pgtransport := &PgTransport{
		stat:  s,
		db:    db,
		query: query,
	}

	return pgtransport, nil
}

func (p *PgTransport) EventToPG(event *riemann.Event) {
	host, domain := event.GetHostDomain()
	// get metric
	var metric sql.NullFloat64
	if m, err := event.GetMetric(); err == nil {
		metric = sql.NullFloat64{m, true}
	}
	// get state
	state := &sql.NullString{String: event.GetState()}
	state.Valid = (state.String != "")
	rows, err := p.db.Query(p.query,
		host,
		domain,
		event.GetService(),
		state,
		event.GetDescription(),
		metric,
		event.GetSQLTags(),
	)
	if err != nil {
		log.Println(err)
	} else {
		rows.Close()
	}
	p.stat.donePostgres()
}

func (p *PgTransport) Send(event *riemann.Event) {
	p.stat.inPostgres()
	if p.stat.getPostgres() > p.stat.MaxQueue {
		log.Printf("Drop message for postgresql. MaxQueue: %v, Queue: %v\n", p.stat.MaxQueue, p.stat.getPostgres())
		return
	}
	p.EventToPG(event)
}
