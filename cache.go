package nucleus

import "sync"

type Cache interface {
	Add(key, value interface{})
	Get(key interface{}) (value interface{}, ok bool)
	Peek(key interface{}) (value interface{}, ok bool)
	Contains(key interface{}) (ok bool)
	Len() int
	Cap() int
	Clear() int
	ReCap(int) error
}

type LruCache struct {
	lru  *lru
	lock sync.RWMutex
}

func NewLruCache(cap int) (*LruCache, error) {
	lru, err := newLruCache(cap)
	if err != nil {
		return nil, err
	}
	cache := &LruCache{
		lru: lru,
	}
	return cache, nil
}

func (c *LruCache) Add(key, value interface{}) {
	c.lock.Lock()
	c.lru.add(key, value)
	c.lock.Unlock()
}

func (c *LruCache) Get(key interface{}) (value interface{}, ok bool) {
	c.lock.Lock()
	value, ok = c.lru.get(key, true)
	c.lock.Unlock()
	return value, ok
}

func (c *LruCache) Peek(key interface{}) (value interface{}, ok bool) {
	c.lock.RLock()
	value, ok = c.lru.get(key, false)
	c.lock.RUnlock()
	return value, ok
}

func (c *LruCache) Contains(key interface{}) (ok bool) {
	c.lock.RLock()
	containKey := c.lru.contains(key)
	c.lock.RUnlock()
	return containKey
}

func (c *LruCache) Len() int {
	c.lock.RLock()
	length := c.lru.len()
	c.lock.RUnlock()
	return length
}

func (c *LruCache) Cap() int {
	return c.lru.cap()
}

func (c *LruCache) Clear() int {
	c.lock.Lock()
	length := c.lru.purge()
	c.lock.Unlock()
	return length
}

func (c *LruCache) ReCap(newCap int) (err error) {
	c.lock.Lock()
	err = c.lru.reCap(newCap)
	c.lock.Unlock()
	return
}
