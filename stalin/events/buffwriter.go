package events

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

type BuffWriterEvent interface {
	ToBytes() ([]byte, error)
}

type BuffWriter struct {
	errChan        chan error
	flushTime      time.Duration
	slReconectTime time.Duration
	writerSize     int
	eventChan      chan BuffWriterEvent
	connBuild      func() (net.Conn, error)
	writer         *bufio.Writer
	conn           net.Conn
	maxQueueSize   int
	queueSize      int
}

func NewBuffWriter(connBuild func() (net.Conn, error), errChan chan error, flush, sleep, size, queue int) *BuffWriter {
	return &BuffWriter{
		connBuild:      connBuild,
		errChan:        errChan,
		eventChan:      make(chan BuffWriterEvent),
		flushTime:      time.Duration(flush) * time.Millisecond,
		writerSize:     size,
		slReconectTime: time.Duration(flush) * time.Millisecond,
		maxQueueSize:   queue,
	}
}

func (b *BuffWriter) loopConnect() {
	for {
		conn, err := b.connBuild()
		if err != nil {
			time.Sleep(b.slReconectTime)
			b.errChan <- err
			continue
		}
		b.conn = conn
		b.writer = bufio.NewWriterSize(b.conn, b.writerSize)
		break
	}
}

func (b *BuffWriter) Inject(e BuffWriterEvent) {
	b.queueSize++
	if b.queueSize > b.maxQueueSize {
		if b.queueSize%100 == 0 {
			b.errChan <- fmt.Errorf("queue size limit")
		}
	} else {
		b.eventChan <- e
	}
	b.queueSize--
}

func (b *BuffWriter) Run() {
	b.loopConnect()
	ticker := time.Tick(b.flushTime)
	for {
		select {
		case <-ticker:
			if err := b.writer.Flush(); err != nil {
				b.errChan <- err
			}
		case event := <-b.eventChan:
			data, err := event.ToBytes()
			if err == nil {
				if _, err := b.writer.Write(data); err != nil {
					b.errChan <- err
					b.loopConnect()
				}
			}
		}
	}
}
