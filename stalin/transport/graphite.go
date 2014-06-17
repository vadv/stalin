package transport

import (
	"bufio"
	"log"
	"net"
	"stalin/riemann"
	"time"
)

type GraphiteTransport struct {
	w             *bufio.Writer
	e             chan []byte
	Addr          string
	tick          <-chan time.Time
	stat          *Stat
	FlushInterval time.Duration
}

func NewGraphite(addr string, interval time.Duration, s *Stat) (*GraphiteTransport, error) {

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	w := bufio.NewWriterSize(conn, 131072)

	graphite := GraphiteTransport{
		w:             w,
		e:             make(chan []byte, 10000),
		tick:          time.Tick(interval),
		Addr:          addr,
		FlushInterval: interval,
		stat:          s,
	}

	go graphite.ReceiveAndTick()

	return &graphite, nil
}

func (g *GraphiteTransport) ReceiveAndTick() {
	for {
		select {
		case <-g.tick:
			err := g.w.Flush()
			if err != nil {
				log.Println("Error flushing to graphite: ", err)
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

func (g *GraphiteTransport) Send(event *riemann.Event) {
	g.stat.inGraphite()
	if g.stat.getGraphite() > g.stat.MaxQueue {
		log.Printf("Drop message for graphite. MaxQueue: %v, Queue: %v\n", g.stat.MaxQueue, g.stat.getGraphite())
		return
	}
	graphitemsg, err := event.ToGraphite()
	if err == nil {
		g.e <- []byte(graphitemsg)
	}
	g.stat.doneGraphite()
}
