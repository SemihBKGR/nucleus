package nucleus

import "sync"

// Cache is the base interface.
type Cache interface {
	Add(key, value interface{})
	Set(key, value interface{}) (ok bool)
	Get(key interface{}) (value interface{}, ok bool)
	Peek(key interface{}) (value interface{}, ok bool)
	Remove(key interface{}) (ok bool)
	Contains(key interface{}) (ok bool)
	Len() int
	Cap() int
	Clear() int
	ReCap(int) error
	Keys() []interface{}
	Values() []interface{}
}

// LruCache is Cache implementation using LRU algorithm.
type LruCache struct {
	lru  *lru
	lock sync.RWMutex
}

// NewLruCache creates new LruCache.
// *LruCache and error are returned.
// Cap must be positive value.
func NewLruCache(cap int) (*LruCache, error) {
	lru, err := newLru(cap)
	if err != nil {
		return nil, err
	}
	cache := &LruCache{
		lru: lru,
	}
	return cache, nil
}

// Add caches key and value.
func (c *LruCache) Add(key, value interface{}) {
	c.lock.Lock()
	c.lru.add(key, value)
	c.lock.Unlock()
}

// Set updates cache entry.
// Returns true if value updated.
func (c *LruCache) Set(key, value interface{}) (ok bool) {
	ok = c.Contains(key)
	if ok {
		c.Add(key, value)
	}
	return

}

// Get returns value of cached entry.
func (c *LruCache) Get(key interface{}) (value interface{}, ok bool) {
	c.lock.Lock()
	value, ok = c.lru.get(key, true)
	c.lock.Unlock()
	return value, ok
}

// Peek return value of cache entry without.
func (c *LruCache) Peek(key interface{}) (value interface{}, ok bool) {
	c.lock.RLock()
	value, ok = c.lru.get(key, false)
	c.lock.RUnlock()
	return value, ok
}

// Remove removes cache entry.
func (c *LruCache) Remove(key interface{}) (ok bool) {
	c.lock.RLock()
	ok = c.lru.remove(key)
	c.lock.RUnlock()
	return
}

// Contains returns true if there is a cache entry given given key.
func (c *LruCache) Contains(key interface{}) (ok bool) {
	c.lock.RLock()
	containKey := c.lru.contains(key)
	c.lock.RUnlock()
	return containKey
}

// Len returns length of the cache.
func (c *LruCache) Len() int {
	c.lock.RLock()
	length := c.lru.len()
	c.lock.RUnlock()
	return length
}

// Cap returns capacity of the cache.
func (c *LruCache) Cap() int {
	return c.lru.cap()
}

// Clear removes all entries in the cache.
func (c *LruCache) Clear() int {
	c.lock.Lock()
	length := c.lru.purge()
	c.lock.Unlock()
	return length
}

// ReCap set capacity of the cache.
// Returns error unless newCap is negative value.
func (c *LruCache) ReCap(newCap int) (err error) {
	c.lock.Lock()
	err = c.lru.reCap(newCap)
	c.lock.Unlock()
	return
}

// Keys returns a slice of entry keys in the cache.
func (c *LruCache) Keys() []interface{} {
	c.lock.RLock()
	keys := c.lru.keys()
	c.lock.RUnlock()
	return keys
}

// Values returns a slice of entry values in the cache.
func (c *LruCache) Values() []interface{} {
	c.lock.RLock()
	values := c.lru.values()
	c.lock.RUnlock()
	return values
}
