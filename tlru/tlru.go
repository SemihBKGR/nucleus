package tlru

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

// Tlru Time Aware Least recently used cache policy
type Tlru struct {
	capacity           int
	elementMap         map[interface{}]*list.Element
	evictionList       *list.List
	expirationDuration time.Duration
	daemonStarted      bool
}

type entry struct {
	key    interface{}
	value  interface{}
	timeMs int64
}

// NewTlru returns new tlru
func NewTlru(capacity int, expiration time.Duration) (*Tlru, error) {
	if capacity <= 0 {
		return nil, errors.New("capacity must be positive value")
	}
	tlru := &Tlru{
		capacity:           capacity,
		elementMap:         make(map[interface{}]*list.Element),
		evictionList:       list.New(),
		expirationDuration: expiration,
		daemonStarted:      false,
	}
	return tlru, nil
}

// StartDaemon starts time expiration daemon
func (t *Tlru) StartDaemon(lock *sync.RWMutex) (ok bool) {
	if t.daemonStarted {
		return false
	}
	if t.expirationDuration <= 0 {
		return false
	}
	t.daemonStarted = true
	go func() {
		for {
			time.Sleep(t.expirationDuration)
			expiredKeys := make([]interface{}, 0)
			lock.RLock()
			currentTimeMs := time.Now().UnixMilli()
			element := t.evictionList.Back()
			for element != nil {
				entry := element.Value.(*entry)
				if entry.timeMs+t.expirationDuration.Milliseconds() <= currentTimeMs {
					expiredKeys = append(expiredKeys, entry.key)
				}
				element = element.Prev()
			}
			lock.RUnlock()
			lock.Lock()
			for _, k := range expiredKeys {
				t.Remove(k)
			}
			lock.Unlock()
		}
	}()
	return true
}

// Add adds entry in cache
func (t *Tlru) Add(key, value interface{}) (eviction bool) {
	if element, ok := t.elementMap[key]; ok {
		t.evictionList.MoveToFront(element)
		element.Value.(*entry).value = value
		return false
	}
	eviction = len(t.elementMap) >= t.capacity
	if eviction {
		t.evict()
	}
	entry := &entry{
		key:    key,
		value:  value,
		timeMs: time.Now().UnixMilli(),
	}
	element := t.evictionList.PushFront(entry)
	t.elementMap[key] = element
	return
}

// Get returns value of cached entry.
func (t *Tlru) Get(key interface{}, trigger bool) (value interface{}, ok bool) {
	if element, ok := t.elementMap[key]; ok {
		if trigger {
			t.evictionList.MoveToFront(element)
		}
		return element.Value.(*entry).value, true
	}
	return nil, false
}

// Remove removes cache entry.
func (t *Tlru) Remove(key interface{}) bool {
	element, ok := t.elementMap[key]
	if !ok {
		return false
	}
	delete(t.elementMap, key)
	t.evictionList.Remove(element)
	return true
}

func (t *Tlru) evict() bool {
	element := t.evictionList.Back()
	if element != nil {
		return t.Remove(element.Value.(*entry).key)
	}
	return false
}

// Clear removes all entries in the cache.
func (t *Tlru) Clear() int {
	length := t.Len()
	for key := range t.elementMap {
		delete(t.elementMap, key)
	}
	t.evictionList.Init()
	return length
}

// Len returns length of the cache.
func (t *Tlru) Len() int {
	return t.evictionList.Len()
}

// Cap returns capacity of the cache.
func (t *Tlru) Cap() int {
	return t.capacity
}

// SetCap set capacity of the cache.
// Returns error unless newCap is negative value.
func (t *Tlru) SetCap(newCapacity int) error {
	if newCapacity <= 0 {
		return errors.New("capacity must be positive value")
	}
	if t.capacity < newCapacity {
		diff := newCapacity - t.capacity
		for i := 0; i < diff; i++ {
			t.evict()
		}
	}
	t.capacity = newCapacity
	return nil
}

// Keys returns a slice of entry keys in the cache.
func (t *Tlru) Keys() []interface{} {
	keys := make([]interface{}, 0, len(t.elementMap))
	for k := range t.elementMap {
		keys = append(keys, k)
	}
	return keys
}

// Values returns a slice of entry values in the cache.
func (t *Tlru) Values() []interface{} {
	values := make([]interface{}, 0, len(t.elementMap))
	for _, v := range t.elementMap {
		values = append(values, v.Value.(*entry).value)
	}
	return values
}

// DaemonStarted returns true if expiration eviction daemon started
func (t *Tlru) DaemonStarted() bool {
	return t.daemonStarted
}

// ExpirationDuration returns duration of expiration
func (t *Tlru) ExpirationDuration() time.Duration {
	return t.expirationDuration
}
