package lru

import (
	"container/list"
	"errors"
)

// Lru Least recently used cache policy
type Lru struct {
	capacity     int
	elementMap   map[interface{}]*list.Element
	evictionList *list.List
}

type entry struct {
	key   interface{}
	value interface{}
}

// NewLru returns new lru
func NewLru(capacity int) (*Lru, error) {
	if capacity <= 0 {
		return nil, errors.New("capacity must be positive value")
	}
	lru := &Lru{
		capacity:     capacity,
		elementMap:   make(map[interface{}]*list.Element),
		evictionList: list.New(),
	}
	return lru, nil
}

// Add adds entry in cache
func (l *Lru) Add(key, value interface{}) (eviction bool) {
	if element, ok := l.elementMap[key]; ok {
		l.evictionList.MoveToFront(element)
		element.Value.(*entry).value = value
		return false
	}
	eviction = len(l.elementMap) >= l.capacity
	if eviction {
		l.evict()
	}
	entry := &entry{
		key:   key,
		value: value,
	}
	element := l.evictionList.PushFront(entry)
	l.elementMap[key] = element
	return
}

// Get returns value of cached entry.
func (l *Lru) Get(key interface{}, trigger bool) (value interface{}, ok bool) {
	if element, ok := l.elementMap[key]; ok {
		if trigger {
			l.evictionList.MoveToFront(element)
		}
		return element.Value.(*entry).value, true
	}
	return nil, false
}

// Remove removes cache entry.
func (l *Lru) Remove(key interface{}) bool {
	element, ok := l.elementMap[key]
	if !ok {
		return false
	}
	delete(l.elementMap, key)
	l.evictionList.Remove(element)
	return true
}

func (l *Lru) evict() bool {
	element := l.evictionList.Back()
	if element != nil {
		return l.Remove(element.Value.(*entry).key)
	}
	return false
}

// Clear removes all entries in the cache.
func (l *Lru) Clear() int {
	length := l.Len()
	for key := range l.elementMap {
		delete(l.elementMap, key)
	}
	l.evictionList.Init()
	return length
}

// Len returns length of the cache.
func (l *Lru) Len() int {
	return l.evictionList.Len()
}

// Cap returns capacity of the cache.
func (l *Lru) Cap() int {
	return l.capacity
}

// SetCap set capacity of the cache.
// Returns error unless newCap is negative value.
func (l *Lru) SetCap(newCapacity int) error {
	if newCapacity <= 0 {
		return errors.New("capacity must be positive value")
	}
	if l.capacity < newCapacity {
		diff := newCapacity - l.capacity
		for i := 0; i < diff; i++ {
			l.evict()
		}
	}
	l.capacity = newCapacity
	return nil
}

// Keys returns a slice of entry keys in the cache.
func (l *Lru) Keys() []interface{} {
	keys := make([]interface{}, 0, len(l.elementMap))
	for k := range l.elementMap {
		keys = append(keys, k)
	}
	return keys
}

// Values returns a slice of entry values in the cache.
func (l *Lru) Values() []interface{} {
	values := make([]interface{}, 0, len(l.elementMap))
	for _, v := range l.elementMap {
		values = append(values, v.Value.(*entry).value)
	}
	return values
}
