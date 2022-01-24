package fifo

import (
	"container/list"
	"errors"
)

// Fifo First in first out cache policy
type Fifo struct {
	capacity     int
	elementMap   map[interface{}]*list.Element
	evictionList *list.List
}

type entry struct {
	key   interface{}
	value interface{}
}

// NewFifo returns new lru
func NewFifo(capacity int) (*Fifo, error) {
	if capacity <= 0 {
		return nil, errors.New("capacity must be positive value")
	}
	fifo := &Fifo{
		capacity:     capacity,
		elementMap:   make(map[interface{}]*list.Element),
		evictionList: list.New(),
	}
	return fifo, nil
}

// Add adds entry in cache
func (f *Fifo) Add(key, value interface{}) (eviction bool) {
	if element, ok := f.elementMap[key]; ok {
		f.evictionList.MoveToFront(element)
		element.Value.(*entry).value = value
		return false
	}
	eviction = len(f.elementMap) >= f.capacity
	if eviction {
		f.evict()
	}
	entry := &entry{
		key:   key,
		value: value,
	}
	element := f.evictionList.PushFront(entry)
	f.elementMap[key] = element
	return
}

// Get returns value of cached entry.
func (f *Fifo) Get(key interface{}, _ bool) (value interface{}, ok bool) {
	if element, ok := f.elementMap[key]; ok {
		return element.Value.(*entry).value, true
	}
	return nil, false
}

// Remove removes cache entry.
func (f *Fifo) Remove(key interface{}) bool {
	element, ok := f.elementMap[key]
	if !ok {
		return false
	}
	delete(f.elementMap, key)
	f.evictionList.Remove(element)
	return true
}

func (f *Fifo) evict() bool {
	element := f.evictionList.Back()
	if element != nil {
		return f.Remove(element.Value.(*entry).key)
	}
	return false
}

// Clear removes all entries in the cache.
func (f *Fifo) Clear() int {
	length := f.Len()
	for key := range f.elementMap {
		delete(f.elementMap, key)
	}
	f.evictionList.Init()
	return length
}

// Len returns length of the cache.
func (f *Fifo) Len() int {
	return f.evictionList.Len()
}

// Cap returns capacity of the cache.
func (f *Fifo) Cap() int {
	return f.capacity
}

// SetCap set capacity of the cache.
// Returns error unless newCap is negative value.
func (f *Fifo) SetCap(newCapacity int) error {
	if newCapacity <= 0 {
		return errors.New("capacity must be positive value")
	}
	if f.capacity < newCapacity {
		diff := newCapacity - f.capacity
		for i := 0; i < diff; i++ {
			f.evict()
		}
	}
	f.capacity = newCapacity
	return nil
}

// Keys returns a slice of entry keys in the cache.
func (f *Fifo) Keys() []interface{} {
	keys := make([]interface{}, 0, len(f.elementMap))
	for k := range f.elementMap {
		keys = append(keys, k)
	}
	return keys
}

// Values returns a slice of entry values in the cache.
func (f *Fifo) Values() []interface{} {
	values := make([]interface{}, 0, len(f.elementMap))
	for _, v := range f.elementMap {
		values = append(values, v.Value.(*entry).value)
	}
	return values
}
