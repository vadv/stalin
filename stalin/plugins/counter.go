package plugins

import (
	"sync"
)

type Counter struct {
	counter int64
	mutex   *sync.RWMutex
}

func NewCounter() *Counter {
	return &Counter{mutex: &sync.RWMutex{}}
}

func (c *Counter) Inc() {
	c.mutex.RLock()
	c.counter++
	c.mutex.RUnlock()
}

func (c *Counter) Dec() {
	c.mutex.RLock()
	c.counter--
	c.mutex.RUnlock()
}

func (c *Counter) Count() int64 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.counter
}

func (c *Counter) Clear() {
	c.mutex.RLock()
	c.counter = 0
	c.mutex.RUnlock()
}
