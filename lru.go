package nucleus

import (
	"container/list"
	"errors"
)

type lru interface {
	add(key, value interface{})
	get(key interface{}, moveFront bool) (value interface{}, ok bool)
	remove() (key, value interface{}, ok bool)
	contains(key interface{}) (ok bool)
	len() int
	cap() int
}

type lruCache struct {
	capacity   int
	elementMap map[interface{}]*list.Element
	useList    *list.List
}

type entry struct {
	key   interface{}
	value interface{}
}

func newLruCache(capacity int) (*lruCache, error) {
	if capacity <= 0 {
		return nil, errors.New("capacity must be positive value")
	}
	lru := &lruCache{
		capacity:   capacity,
		elementMap: make(map[interface{}]*list.Element),
		useList:    list.New(),
	}
	return lru, nil
}

func (lru *lruCache) add(key, value interface{}) {
	if len(lru.elementMap) >= lru.capacity {
		lru.remove()
	}
	if element, ok := lru.elementMap[key]; ok {
		lru.useList.MoveToFront(element)
		element.Value.(*entry).value = value
	}
	entry := &entry{
		key:   key,
		value: value,
	}
	element := lru.useList.PushFront(entry)
	lru.elementMap[key] = element
}

func (lru *lruCache) get(key interface{}, moveFront bool) (value interface{}, ok bool) {
	if element, ok := lru.elementMap[key]; ok {
		if moveFront {
			lru.useList.MoveToFront(element)
		}
		return element.Value.(*entry).value, true
	}
	return nil, false
}

func (lru *lruCache) contains(key interface{}) (ok bool) {
	_, ok = lru.elementMap[key]
	return
}

func (lru *lruCache) len() int {
	return lru.useList.Len()
}

func (lru *lruCache) cap() int {
	return lru.capacity
}

func (lru *lruCache) remove() (key, value interface{}, ok bool) {
	element := lru.useList.Back()
	if element != nil {
		entry := element.Value.(*entry)
		lru.useList.Remove(element)
		delete(lru.elementMap, entry.key)
		return entry.key, entry.value, true
	}
	return nil, nil, false
}
