package mru

import (
	"container/list"
	"errors"
)

// Mru Most recently used policy
type Mru struct {
	capacity     int
	elementMap   map[interface{}]*list.Element
	evictionList *list.List
}

type entry struct {
	key   interface{}
	value interface{}
}

// NewMru returns new mru
func NewMru(capacity int) (*Mru, error) {
	if capacity <= 0 {
		return nil, errors.New("capacity must be positive value")
	}
	lru := &Mru{
		capacity:     capacity,
		elementMap:   make(map[interface{}]*list.Element),
		evictionList: list.New(),
	}
	return lru, nil
}

// Add adds entry in cache
func (m *Mru) Add(key, value interface{}) (eviction bool) {
	if element, ok := m.elementMap[key]; ok {
		m.evictionList.MoveToFront(element)
		element.Value.(*entry).value = value
		return false
	}
	eviction = len(m.elementMap) >= m.capacity
	if eviction {
		m.evict()
	}
	entry := &entry{
		key:   key,
		value: value,
	}
	element := m.evictionList.PushFront(entry)
	m.elementMap[key] = element
	return
}

// Get returns value of cached entry.
func (m *Mru) Get(key interface{}, trigger bool) (value interface{}, ok bool) {
	if element, ok := m.elementMap[key]; ok {
		if trigger {
			m.evictionList.MoveToFront(element)
		}
		return element.Value.(*entry).value, true
	}
	return nil, false
}

// Remove removes cache entry.
func (m *Mru) Remove(key interface{}) bool {
	element, ok := m.elementMap[key]
	if !ok {
		return false
	}
	delete(m.elementMap, key)
	m.evictionList.Remove(element)
	return true
}

func (m *Mru) evict() bool {
	element := m.evictionList.Front()
	if element != nil {
		return m.Remove(element.Value.(*entry).key)
	}
	return false
}

// Clear removes all entries in the cache.
func (m *Mru) Clear() int {
	length := m.Len()
	for key := range m.elementMap {
		delete(m.elementMap, key)
	}
	m.evictionList.Init()
	return length
}

// Len returns length of the cache.
func (m *Mru) Len() int {
	return m.evictionList.Len()
}

// Cap returns capacity of the cache.
func (m *Mru) Cap() int {
	return m.capacity
}

// SetCap set capacity of the cache.
// Returns error unless newCap is negative value.
func (m *Mru) SetCap(newCapacity int) error {
	if newCapacity <= 0 {
		return errors.New("capacity must be positive value")
	}
	if m.capacity < newCapacity {
		diff := newCapacity - m.capacity
		for i := 0; i < diff; i++ {
			m.evict()
		}
	}
	m.capacity = newCapacity
	return nil
}

// Keys returns a slice of entry keys in the cache.
func (m *Mru) Keys() []interface{} {
	keys := make([]interface{}, 0, len(m.elementMap))
	for k := range m.elementMap {
		keys = append(keys, k)
	}
	return keys
}

// Values returns a slice of entry values in the cache.
func (m *Mru) Values() []interface{} {
	values := make([]interface{}, 0, len(m.elementMap))
	for _, v := range m.elementMap {
		values = append(values, v.Value.(*entry).value)
	}
	return values
}
