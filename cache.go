package nucleus

import (
	"github.com/SemihBKGR/nucleus/fifo"
	"github.com/SemihBKGR/nucleus/lru"
	"github.com/SemihBKGR/nucleus/mru"
	"github.com/SemihBKGR/nucleus/tlru"
	"sync"
	"time"
)

// Policy base interface of policies.
type Policy interface {
	Add(key, value interface{}) (eviction bool)
	Get(key interface{}, trigger bool) (value interface{}, ok bool)
	Remove(key interface{}) (ok bool)
	Clear() int
	Len() int
	Cap() int
	SetCap(int) error
	Keys() []interface{}
	Values() []interface{}
}

// Cache is main struct.
type Cache struct {
	policy Policy
	lock   sync.RWMutex
}

// NewLruCache returns new cache with lru policy.
func NewLruCache(cap int) (*Cache, error) {
	lruPolicy, err := lru.NewLru(cap)
	if err != nil {
		return nil, err
	}
	cache := &Cache{
		policy: lruPolicy,
	}
	return cache, nil
}

// NewMruCache returns new cache with mru policy.
func NewMruCache(cap int) (*Cache, error) {
	mruPolicy, err := mru.NewMru(cap)
	if err != nil {
		return nil, err
	}
	cache := &Cache{
		policy: mruPolicy,
	}
	return cache, nil
}

// NewFifoCache returns new cache with fifo policy.
func NewFifoCache(cap int) (*Cache, error) {
	fifoPolicy, err := fifo.NewFifo(cap)
	if err != nil {
		return nil, err
	}
	cache := &Cache{
		policy: fifoPolicy,
	}
	return cache, nil
}

// NewTlruCache create new cache with tlru policy
func NewTlruCache(cap int, expDur time.Duration) (*Cache, error) {
	tlruPolicy, err := tlru.NewTlru(cap, expDur)
	if err != nil {
		return nil, err
	}
	cache := &Cache{
		policy: tlruPolicy,
	}
	tlruPolicy.StartEvictionDaemon(&cache.lock)
	return cache, nil
}

// Add adds entry in cache
func (c *Cache) Add(key, value interface{}) (eviction bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.policy.Add(key, value)
}

// Set updates cache entry.
// Returns true if value updated.
func (c *Cache) Set(key, value interface{}) (ok bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, ok = c.policy.Get(key, false)
	if ok {
		c.policy.Add(key, value)
	}
	return
}

// Get returns value of cached entry.
func (c *Cache) Get(key interface{}) (value interface{}, ok bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	value, ok = c.policy.Get(key, true)
	return
}

// Remove removes cache entry.
func (c *Cache) Remove(key interface{}) (ok bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	ok = c.policy.Remove(key)
	return
}

// Contains returns true if there is a cache entry given given key.
func (c *Cache) Contains(key interface{}) (ok bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	_, ok = c.policy.Get(key, false)
	return
}

// Clear removes all entries in the cache.
func (c *Cache) Clear() (length int) {
	c.lock.Lock()
	defer c.lock.Unlock()
	length = c.policy.Clear()
	return
}

// Len returns length of the cache.
func (c *Cache) Len() (len int) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	len = c.policy.Len()
	return
}

// Cap returns capacity of the cache.
func (c *Cache) Cap() int {
	return c.policy.Cap()
}

// SetCap set capacity of the cache.
// Returns error unless newCap is negative value.
func (c *Cache) SetCap(newCap int) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.policy.SetCap(newCap)
}

// Keys returns a slice of entry keys in the cache.
func (c *Cache) Keys() []interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.policy.Keys()
}

// Values returns a slice of entry values in the cache.
func (c *Cache) Values() []interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.policy.Values()
}
