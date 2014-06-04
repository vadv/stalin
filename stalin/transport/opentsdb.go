package transport

import (
	"bufio"
	"log"
	"net"
	"stalin/riemann"
	"time"
)

type OpentsDBTransport struct {
	w             *bufio.Writer
	e             chan []byte
	Addr          string
	tick          <-chan time.Time
	stat          *Stat
	FlushInterval time.Duration
}

func NewOpentsDB(addr string, interval time.Duration, s *Stat) (*OpentsDBTransport, error) {

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	w := bufio.NewWriterSize(conn, 131072)

	opentsdb := OpentsDBTransport{
		w:             w,
		e:             make(chan []byte, 10000),
		tick:          time.Tick(interval),
		Addr:          addr,
		FlushInterval: interval,
		stat:          s,
	}

	go opentsdb.ReceiveAndTick()

	return &opentsdb, nil
}

func (g *OpentsDBTransport) ReceiveAndTick() {
	for {
		select {
		case <-g.tick:
			err := g.w.Flush()
			if err != nil {
				log.Println("Error flushing to opentsdb: ", err)
			}

		case data := <-g.e:
			count, err := g.w.Write(data)
			// reconnect
			if err != nil {
				log.Printf("Failed to send data to %v size: %v, error: %v\n", g.Addr, count, err)
				for {
					conn, err := net.Dial("tcp", g.Addr)
					if err != nil {
						log.Printf("Error reconnecting to %v, error: %v\n", g.Addr, err)
						time.Sleep(time.Second)
						continue
					}
					g.w = bufio.NewWriterSize(conn, 65536)
					break
				}
			}
		}
	}
}

func (g *OpentsDBTransport) Send(event *riemann.Event) {
	g.stat.inOpentsDB()
	opentsdbmsg, err := event.ToOpentsDB()
	if err == nil {
		g.e <- []byte(opentsdbmsg)
	}
	g.stat.doneOpentsDB()
}
