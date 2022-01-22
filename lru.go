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

func newLru(capacity int) (*lru, error) {
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
		l.removeLast()
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

func (l *lru) get(key interface{}, move bool) (value interface{}, ok bool) {
	if element, ok := l.elementMap[key]; ok {
		if move {
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

func (l *lru) remove(key interface{}) bool {
	element, ok := l.elementMap[key]
	if !ok {
		return false
	}
	delete(l.elementMap, key)
	l.useList.Remove(element)
	return true
}

func (l *lru) removeLast() bool {
	element := l.useList.Back()
	if element != nil {
		return l.remove(element.Value.(*entry).key)
	}
	return false
}

func (l *lru) reCap(newCapacity int) error {
	if newCapacity <= 0 {
		return errors.New("capacity must be positive value")
	}
	if l.capacity < newCapacity {
		diff := newCapacity - l.capacity
		for i := 0; i < diff; i++ {
			l.removeLast()
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

func (l *lru) keys() []interface{} {
	keys := make([]interface{}, 0, len(l.elementMap))
	for k := range l.elementMap {
		keys = append(keys, k)
	}
	return keys
}

func (l *lru) values() []interface{} {
	values := make([]interface{}, 0, len(l.elementMap))
	for _, v := range l.elementMap {
		values = append(values, v.Value.(*entry).value)
	}
	return values
}
