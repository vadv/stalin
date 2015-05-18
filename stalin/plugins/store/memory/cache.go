package storage

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"

	. "stalin/plugins/store/query_language"
)

type Cache struct {
	*cache
}

type cache struct {
	sync.RWMutex
	problems *Problems
	janitor  *janitor
}

func (c *cache) Set(k string, x *Problem, d time.Duration) {
	c.Lock()
	c.set(k, x, d)
	c.Unlock()
}

func (c *cache) set(k string, x *Problem, d time.Duration) {
	c.problems.Items[k] = x
}

func (c *cache) Add(k string, x *Problem, d time.Duration) error {
	c.Lock()
	_, found := c.get(k)
	if found {
		c.Unlock()
		return fmt.Errorf("Item %s already exists", k)
	}
	c.set(k, x, d)
	c.Unlock()
	return nil
}

func (c *cache) Replace(k string, x *Problem, d time.Duration) error {
	c.Lock()
	_, found := c.get(k)
	if !found {
		c.Unlock()
		return fmt.Errorf("Item %s doesn't exist", k)
	}
	c.set(k, x, d)
	c.Unlock()
	return nil
}

func (c *cache) Get(k string) (*Problem, bool) {
	c.RLock()
	x, found := c.get(k)
	c.RUnlock()
	return x, found
}

func (c *cache) get(k string) (*Problem, bool) {
	item, found := c.problems.Items[k]
	if !found || item.Expired() {
		return nil, false
	}
	return item, true
}

func (c *cache) Delete(k string) {
	c.Lock()
	c.delete(k)
	c.Unlock()
}

func (c *cache) delete(k string) {
	delete(c.problems.Items, k)
}

func (c *cache) DeleteExpired() {
	c.Lock()
	for k, v := range c.problems.Items {
		if v.Expired() {
			c.delete(k)
		}
	}
	c.Unlock()
}

func (c *cache) Save(w io.Writer) (err error) {
	enc := gob.NewEncoder(w)
	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("Error registering item types with Gob library")
		}
	}()
	c.RLock()
	defer c.RUnlock()
	gob.Register(&Problem{})
	gob.Register(&Problems{})
	err = enc.Encode(&c.problems)
	return
}

func (c *cache) SaveFile(fname string) error {
	fp, err := os.Create(fname)
	if err != nil {
		return err
	}
	err = c.Save(fp)
	if err != nil {
		fp.Close()
		return err
	}
	return fp.Close()
}

func (c *cache) Load(r io.Reader) error {
	dec := gob.NewDecoder(r)
	problems := NewProblems()
	err := dec.Decode(&problems)
	if err == nil {
		c.Lock()
		defer c.Unlock()
		for k, v := range problems.Items {
			ov, found := c.problems.Items[k]
			if !found || ov.Expired() {
				c.problems.Items[k] = v
			}
		}
	}
	return err
}

func (c *cache) LoadFile(fname string) error {
	fp, err := os.Open(fname)
	if err != nil {
		return err
	}
	err = c.Load(fp)
	if err != nil {
		fp.Close()
		return err
	}
	return fp.Close()
}

func (c *cache) List() *Problems {
	return c.problems
}

func (c *cache) ItemCount() int {
	c.RLock()
	n := len(c.problems.Items)
	c.RUnlock()
	return n
}

func (c *cache) Flush() {
	c.Lock()
	c.problems = NewProblems()
	c.Unlock()
}

type janitor struct {
	Interval time.Duration
	stop     chan bool
}

func (j *janitor) Run(c *cache) {
	j.stop = make(chan bool)
	tick := time.Tick(j.Interval)
	for {
		select {
		case <-tick:
			c.DeleteExpired()
		case <-j.stop:
			return
		}
	}
}

func stopJanitor(c *Cache) {
	c.janitor.stop <- true
}

func runJanitor(c *cache, ci time.Duration) {
	j := &janitor{
		Interval: ci,
	}
	c.janitor = j
	go j.Run(c)
}

func newCache(m *Problems) *cache {
	c := &cache{
		problems: m,
	}
	return c
}

func newCacheWithJanitor(ci time.Duration, m *Problems) *Cache {
	c := newCache(m)
	C := &Cache{c}
	if ci > 0 {
		runJanitor(c, ci)
		runtime.SetFinalizer(C, stopJanitor)
	}
	return C
}

func NewCache(cleanupInterval time.Duration) *Cache {
	problems := NewProblems()
	return newCacheWithJanitor(cleanupInterval, problems)
}

func NewCacheFrom(cleanupInterval time.Duration, problems *Problems) *Cache {
	return newCacheWithJanitor(cleanupInterval, problems)
}
