package nucleus

import (
	"container/list"
	"errors"
)

type lru struct {
	capacity   int
	elementMap map[interface{}]*list.Element
	useList    *list.List
}

type entry struct {
	key   interface{}
	value interface{}
}

func newLruCache(capacity int) (*lru, error) {
	if capacity <= 0 {
		return nil, errors.New("capacity must be positive value")
	}
	lru := &lru{
		capacity:   capacity,
		elementMap: make(map[interface{}]*list.Element),
		useList:    list.New(),
	}
	return lru, nil
}

func (l *lru) add(key, value interface{}) {
	if len(l.elementMap) >= l.capacity {
		l.remove()
	}
	if element, ok := l.elementMap[key]; ok {
		l.useList.MoveToFront(element)
		element.Value.(*entry).value = value
	}
	entry := &entry{
		key:   key,
		value: value,
	}
	element := l.useList.PushFront(entry)
	l.elementMap[key] = element
}

func (l *lru) get(key interface{}, moveFront bool) (value interface{}, ok bool) {
	if element, ok := l.elementMap[key]; ok {
		if moveFront {
			l.useList.MoveToFront(element)
		}
		return element.Value.(*entry).value, true
	}
	return nil, false
}

func (l *lru) contains(key interface{}) (ok bool) {
	_, ok = l.elementMap[key]
	return
}

func (l *lru) len() int {
	return l.useList.Len()
}

func (l *lru) cap() int {
	return l.capacity
}

func (l *lru) remove() (key, value interface{}, ok bool) {
	element := l.useList.Back()
	if element != nil {
		entry := element.Value.(*entry)
		l.useList.Remove(element)
		delete(l.elementMap, entry.key)
		return entry.key, entry.value, true
	}
	return nil, nil, false
}

func (l *lru) reCap(newCapacity int) error {
	if newCapacity <= 0 {
		return errors.New("capacity must be positive value")
	}
	if l.capacity < newCapacity {
		diff := newCapacity - l.capacity
		for i := 0; i < diff; i++ {
			l.remove()
		}

	}
	l.capacity = newCapacity
	return nil
}

func (l *lru) purge() int {
	length := l.len()
	for key := range l.elementMap {
		delete(l.elementMap, key)
	}
	l.useList.Init()
	return length
}
