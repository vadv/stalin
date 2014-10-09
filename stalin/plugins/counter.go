package plugins

import (
	"log"
	"sync"
	"time"
)

type Counter struct {
	Name           string
	Done           int
	Queue          int
	Droped         int
	StatTime       int
	PrintStatistic bool
	mutex          sync.RWMutex
}

func NewCounter(name string, stattime int, printstat bool) *Counter {
	c := &Counter{StatTime: stattime, Name: name, PrintStatistic: printstat}
	return c
}

func (c *Counter) In() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Queue++
}

func (c *Counter) Out() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Done++
	c.Queue--
}

func (c *Counter) QueueSize() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.Queue
}

func (c *Counter) UpDropped() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Droped++
}

func (c *Counter) clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Done = 0
	c.Droped = 0
}

func (c *Counter) print() {
	log.Printf("[%v] Input: [%6.2f], Queue: [%6.2f], Drop: [%6.2f]",
		c.Name, float64(c.Done)/float64(c.StatTime), float64(c.Queue)/float64(c.StatTime), float64(c.Droped)/float64(c.StatTime))
}

func (c *Counter) runTicker() {
	ticker := time.Tick(time.Duration(c.StatTime) * time.Second)
	for {
		select {
		case <-ticker:
			if c.PrintStatistic {
				c.print()
			}
			c.clear()
		}
	}
}

func (c *Counter) Start() {
	go c.runTicker()
}
