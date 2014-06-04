package transport

import (
	"code.google.com/p/goprotobuf/proto"
	"errors"
	"log"
	"net"
	"stalin/riemann"
	"time"
)

type Transport struct {
	Graphite    *GraphiteTransport
	OpentsDB    *OpentsDBTransport
	Resender    *ResenderTransport
	Pg          *PgTransport
	useGraphite bool
	useOpentsDB bool
	useResender bool
	usePg       bool
	isSender    bool
	stat        *Stat
}

func NewTransport(graphiteaddr, opentsdbaddr, resenderaddr, pgconstring, pgquery string, pgpool int, graphiteinterval, opentsdbinterval, resenderinterval time.Duration, statistictime int) (*Transport, error) {

	s := &Stat{StatTime: statistictime}
	t := &Transport{stat: s}

	err := errors.New("")

	if graphiteaddr != "" || resenderaddr != "" || pgconstring != "" {
		t.isSender = true
	}

	if graphiteaddr != "" {
		t.Graphite, err = NewGraphite(graphiteaddr, graphiteinterval, s)
		if err != nil {
			return nil, err
		}
		t.useGraphite = true
	}

	if opentsdbaddr != "" {
		t.OpentsDB, err = NewOpentsDB(opentsdbaddr, opentsdbinterval, s)
		if err != nil {
			return nil, err
		}
		t.useOpentsDB = true
	}

	if resenderaddr != "" {
		t.Resender, err = NewResender(resenderaddr, resenderinterval, s)
		if err != nil {
			return nil, err
		}
		t.useResender = true
	}

	if pgconstring != "" {
		t.Pg, err = NewPG(pgconstring, pgquery, pgpool, s)
		if err != nil {
			return nil, err
		}
		t.usePg = true
	}

	go s.Ticks()

	return t, nil

}

func (t *Transport) HandleConn(conn net.Conn) {
	log.Println("Starting handler for", conn.RemoteAddr())
	for {
		msg := &riemann.Msg{}
		t.stat.inInput()
		data, err := riemann.Unpack(conn)
		if err != nil {
			log.Printf("Error decoding protobuf message from client: %v error: %v\n", conn.RemoteAddr(), err.Error())
			t.stat.doneInput()
			return
		}
		err = proto.Unmarshal(data, msg)
		if err != nil {
			log.Printf("Error decoding protobuf message from client: %v error: %v\n", conn.RemoteAddr(), err.Error())
			t.stat.doneInput()
			return
		}
		t.Send(msg)
		t.stat.doneInput()
	}
}

func (t *Transport) Send(msg *riemann.Msg) {
	if t.isSender {
		if t.useResender {
			go t.Resender.Send(msg)
		}
		if t.useGraphite || t.usePg || t.useOpentsDB {
			for _, event := range msg.GetEvents() {
				if t.useGraphite {
					go t.Graphite.Send(event)
				}
				if t.usePg {
					go t.Pg.Send(event)
				}
				if t.useOpentsDB {
					go t.OpentsDB.Send(event)
				}
			}
		}
	} else {
		log.Printf("msg:", msg.String())
	}
}
