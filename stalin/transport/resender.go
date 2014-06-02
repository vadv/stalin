package transport

import (
	"bufio"
	"bytes"
	"log"
	"net"
	"stalin/riemann"
	"time"
)

type ResenderTransport struct {
	w             *bufio.Writer
	e             chan []byte
	Addr          string
	FlushInterval time.Duration
	tick          <-chan time.Time
	stat          *Stat
	con           net.Conn
}

func NewResender(addr string, interval time.Duration, s *Stat) (*ResenderTransport, error) {

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	w := bufio.NewWriterSize(conn, 65536)

	resender := ResenderTransport{
		w:             w,
		e:             make(chan []byte, 10000),
		Addr:          addr,
		tick:          time.Tick(interval),
		FlushInterval: interval,
		stat:          s,
	}

	go resender.ReceiveAndTick()

	return &resender, nil
}

func (g *ResenderTransport) ReceiveAndTick() {
	for {
		select {
		case <-g.tick:
			err := g.w.Flush()
			if err != nil {
				log.Println("Error flushing to resender: ", err)
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

func (r *ResenderTransport) Send(msg *riemann.Msg) {
	r.stat.inResend()
	buf := &bytes.Buffer{}
	enc := riemann.NewEncoder(buf)
	enc.EncodeMsg(msg)
	r.e <- buf.Bytes()
	r.stat.doneResend()
}
