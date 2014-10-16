package cache

import (
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"hash/fnv"
	"io"
	"os"
	"runtime"
	. "stalin/plugins/store/problem"
	"sync"
	"time"
)

type unexportedInterface interface {
	Set(string, interface{}, time.Duration)
	Add(string, interface{}, time.Duration) error
	Replace(string, interface{}, time.Duration) error
	Get(string) (*Problem, bool)
	Delete(string)
	DeleteExpired()
	Items() map[string]*Item
	ItemCount() int
	Flush()
	Save(io.Writer) error
	SaveFile(string) error
	Load(io.Reader) error
	LoadFile(io.Reader) error
}

type Item struct {
	ProblemObj *Problem
	Expiration *time.Time
}

// Returns true if the item has expired.
func (item *Item) Expired() bool {
	if item.Expiration == nil {
		return false
	}
	return item.Expiration.Before(time.Now())
}

type Cache struct {
	*cache
	// If this is confusing, see the comment at the bottom of New()
}

type cache struct {
	sync.RWMutex
	defaultExpiration time.Duration
	items             map[string]*Item
	janitor           *janitor
}

// Add an item to the cache, replacing any existing item. If the duration is 0,
// the cache's default expiration time is used. If it is -1, the item never
// expires.
func (c *cache) Set(k string, x *Problem, d time.Duration) {
	c.Lock()
	c.set(k, x, d)
	// TODO: Calls to mu.Unlock are currently not deferred because defer
	// adds ~200 ns (as of go1.)
	c.Unlock()
}

func (c *cache) set(k string, x *Problem, d time.Duration) {
	var e *time.Time
	if d == 0 {
		d = c.defaultExpiration
	}
	if d > 0 {
		t := time.Now().Add(d)
		e = &t
	}
	c.items[k] = &Item{
		ProblemObj: x,
		Expiration: e,
	}
}

// Add an item to the cache only if an item doesn't already exist for the given
// key, or if the existing item has expired. Returns an error otherwise.
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

// Set a new value for the cache key only if it already exists, and the existing
// item hasn't expired. Returns an error otherwise.
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

// Get an item from the cache. Returns the item or nil, and a bool indicating
// whether the key was found.
func (c *cache) Get(k string) (*Problem, bool) {
	c.RLock()
	x, found := c.get(k)
	c.RUnlock()
	return x, found
}

func (c *cache) get(k string) (*Problem, bool) {
	item, found := c.items[k]
	if !found || item.Expired() {
		return nil, false
	}
	return item.ProblemObj, true
}

// Delete an item from the cache. Does nothing if the key is not in the cache.
func (c *cache) Delete(k string) {
	c.Lock()
	c.delete(k)
	c.Unlock()
}

func (c *cache) delete(k string) {
	delete(c.items, k)
}

// Delete all expired items from the cache.
func (c *cache) DeleteExpired() {
	c.Lock()
	for k, v := range c.items {
		if v.Expired() {
			c.delete(k)
		}
	}
	c.Unlock()
}

// Write the cache's items (using Gob) to an io.Writer.
func (c *cache) Save(w io.Writer) (err error) {
	enc := gob.NewEncoder(w)
	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("Error registering item types with Gob library")
		}
	}()
	c.RLock()
	defer c.RUnlock()
	for _, v := range c.items {
		gob.Register(v.ProblemObj)
	}
	err = enc.Encode(&c.items)
	return
}

// Save the cache's items to the given filename, creating the file if it
// doesn't exist, and overwriting it if it does.
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

// Add (Gob-serialized) cache items from an io.Reader, excluding any items with
// keys that already exist (and haven't expired) in the current cache.
func (c *cache) Load(r io.Reader) error {
	dec := gob.NewDecoder(r)
	items := map[string]*Item{}
	err := dec.Decode(&items)
	if err == nil {
		c.Lock()
		defer c.Unlock()
		for k, v := range items {
			ov, found := c.items[k]
			if !found || ov.Expired() {
				c.items[k] = v
			}
		}
	}
	return err
}

// Load and add cache items from the given filename, excluding any items with
// keys that already exist in the current cache.
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

// Returns the items in the cache. This may include items that have expired,
// but have not yet been cleaned up. If this is significant, the Expiration
// fields of the items should be checked. Note that explicit synchronization
// is needed to use a cache and its corresponding Items() return value at
// the same time, as the map is shared.
func (c *cache) Items() map[string]*Item {
	c.RLock()
	defer c.RUnlock()
	return c.items
}

// Returns the number of items in the cache. This may include items that have
// expired, but have not yet been cleaned up. Equivalent to len(c.Items()).
func (c *cache) ItemCount() int {
	c.RLock()
	n := len(c.items)
	c.RUnlock()
	return n
}

// Delete all items from the cache.
func (c *cache) Flush() {
	c.Lock()
	c.items = map[string]*Item{}
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

func newCache(de time.Duration) *cache {
	if de == 0 {
		de = -1
	}
	c := &cache{
		defaultExpiration: de,
		items:             map[string]*Item{},
	}
	return c
}

// Return a new cache with a given default expiration duration and cleanup
// interval. If the expiration duration is less than 1, the items in the cache
// never expire (by default), and must be deleted manually. If the cleanup
// interval is less than one, expired items are not deleted from the cache
// before calling DeleteExpired.
func New(defaultExpiration, cleanupInterval time.Duration) *Cache {
	c := newCache(defaultExpiration)
	// This trick ensures that the janitor goroutine (which--granted it
	// was enabled--is running DeleteExpired on c forever) does not keep
	// the returned C object from being garbage collected. When it is
	// garbage collected, the finalizer stops the janitor goroutine, after
	// which c can be collected.
	C := &Cache{c}
	if cleanupInterval > 0 {
		runJanitor(c, cleanupInterval)
		runtime.SetFinalizer(C, stopJanitor)
	}
	return C
}

type unexportedShardedCache struct {
	*shardedCache
}

type shardedCache struct {
	m       uint32
	cs      []*cache
	janitor *shardedJanitor
}

func (sc *shardedCache) bucket(k string) *cache {
	h := fnv.New32()
	h.Write([]byte(k))
	n := binary.BigEndian.Uint32(h.Sum(nil))
	return sc.cs[n%sc.m]
}

func (sc *shardedCache) Set(k string, x *Problem, d time.Duration) {
	sc.bucket(k).Set(k, x, d)
}

func (sc *shardedCache) Add(k string, x *Problem, d time.Duration) error {
	return sc.bucket(k).Add(k, x, d)
}

func (sc *shardedCache) Replace(k string, x *Problem, d time.Duration) error {
	return sc.bucket(k).Replace(k, x, d)
}

func (sc *shardedCache) Get(k string) (*Problem, bool) {
	return sc.bucket(k).Get(k)
}

func (sc *shardedCache) Delete(k string) {
	sc.bucket(k).Delete(k)
}

func (sc *shardedCache) DeleteExpired() {
	for _, v := range sc.cs {
		v.DeleteExpired()
	}
}

func (sc *shardedCache) Flush() {
	for _, v := range sc.cs {
		v.Flush()
	}
}

type shardedJanitor struct {
	Interval time.Duration
	stop     chan bool
}

func (j *shardedJanitor) Run(sc *shardedCache) {
	j.stop = make(chan bool)
	tick := time.Tick(j.Interval)
	for {
		select {
		case <-tick:
			sc.DeleteExpired()
		case <-j.stop:
			return
		}
	}
}

func stopShardedJanitor(sc *unexportedShardedCache) {
	sc.janitor.stop <- true
}

func runShardedJanitor(sc *shardedCache, ci time.Duration) {
	j := &shardedJanitor{
		Interval: ci,
	}
	sc.janitor = j
	go j.Run(sc)
}

func newShardedCache(n int, de time.Duration) *shardedCache {
	sc := &shardedCache{
		m:  uint32(n - 1),
		cs: make([]*cache, n),
	}
	for i := 0; i < n; i++ {
		c := &cache{
			defaultExpiration: de,
			items:             map[string]*Item{},
		}
		sc.cs[i] = c
	}
	return sc
}

func unexportedNewSharded(shards int, defaultExpiration, cleanupInterval time.Duration) *unexportedShardedCache {
	if defaultExpiration == 0 {
		defaultExpiration = -1
	}
	sc := newShardedCache(shards, defaultExpiration)
	SC := &unexportedShardedCache{sc}
	if cleanupInterval > 0 {
		runShardedJanitor(sc, cleanupInterval)
		runtime.SetFinalizer(SC, stopShardedJanitor)
	}
	return SC
}
