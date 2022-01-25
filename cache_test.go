package nucleus

import (
	"math"
	"strconv"
	"testing"
)

func TestNewLruCache(t *testing.T) {
	cache, err := NewLruCache(1)
	if cache == nil {
		t.FailNow()
	}
	if err != nil {
		t.FailNow()
	}
	cache, err = NewLruCache(0)
	if cache != nil {
		t.FailNow()
	}
	if err == nil {
		t.FailNow()
	}
	cache, err = NewLruCache(-1)
	if cache != nil {
		t.FailNow()
	}
	if err == nil {
		t.FailNow()
	}
}

func TestNewMruCache(t *testing.T) {
	cache, err := NewMruCache(1)
	if cache == nil {
		t.FailNow()
	}
	if err != nil {
		t.FailNow()
	}
	cache, err = NewMruCache(0)
	if cache != nil {
		t.FailNow()
	}
	if err == nil {
		t.FailNow()
	}
	cache, err = NewMruCache(-1)
	if cache != nil {
		t.FailNow()
	}
	if err == nil {
		t.FailNow()
	}
}

func TestNewFifoCache(t *testing.T) {
	cache, err := NewFifoCache(1)
	if cache == nil {
		t.FailNow()
	}
	if err != nil {
		t.FailNow()
	}
	cache, err = NewFifoCache(0)
	if cache != nil {
		t.FailNow()
	}
	if err == nil {
		t.FailNow()
	}
	cache, err = NewFifoCache(-1)
	if cache != nil {
		t.FailNow()
	}
	if err == nil {
		t.FailNow()
	}
}

func TestNewTlruCache(t *testing.T) {
	cache, err := NewTlruCache(1, 0)
	if cache == nil {
		t.FailNow()
	}
	if err != nil {
		t.FailNow()
	}
	cache, err = NewTlruCache(0, 0)
	if cache != nil {
		t.FailNow()
	}
	if err == nil {
		t.FailNow()
	}
	cache, err = NewTlruCache(-1, 0)
	if cache != nil {
		t.FailNow()
	}
	if err == nil {
		t.FailNow()
	}
}

func TestCache_Add(t *testing.T) {
	capacity := 10
	cache, _ := NewLruCache(capacity)
	for i := 0; i < capacity*2; i++ {
		cache.Add(i, strconv.Itoa(i))
	}
}

func TestCache_Set(t *testing.T) {
	capacity := 10
	cache, _ := NewLruCache(capacity)
	for i := 0; i < capacity; i++ {
		ok := cache.Set(i, strconv.Itoa(i))
		if ok {
			t.FailNow()
		}
	}
	for i := 0; i < capacity; i++ {
		cache.Add(i, strconv.Itoa(i))
	}
	for i := 0; i < capacity; i++ {
		ok := cache.Set(i, strconv.Itoa(i))
		if !ok {
			t.FailNow()
		}
	}
}

func TestCache_Get(t *testing.T) {
	capacity := 10
	cache, _ := NewLruCache(capacity)
	for i := 0; i < capacity; i++ {
		cache.Add(i, strconv.Itoa(i))
	}
	for i := 0; i < capacity; i++ {
		value, ok := cache.Get(i)
		if !ok || value.(string) != strconv.Itoa(i) {
			t.FailNow()
		}
	}
	for i := capacity; i < capacity*2; i++ {
		value, ok := cache.Get(i)
		if ok || value != nil {
			t.FailNow()
		}
	}
}

func TestCache_Remove(t *testing.T) {
	capacity := 10
	cache, _ := NewLruCache(capacity)
	for i := 0; i < capacity; i++ {
		cache.Add(i, strconv.Itoa(i))
	}
	for i := 0; i < capacity; i++ {
		ok := cache.Remove(i)
		if !ok {
			t.FailNow()
		}
	}
	for i := 0; i < capacity; i++ {
		ok := cache.Remove(i)
		if ok {
			t.FailNow()
		}
	}
}

func TestCache_Contains(t *testing.T) {
	capacity := 10
	cache, _ := NewLruCache(capacity)
	for i := 0; i < capacity; i++ {
		cache.Add(i, strconv.Itoa(i))
	}
	for i := 0; i < capacity; i++ {
		if !cache.Contains(i) {
			t.FailNow()
		}
	}
	for i := capacity; i < capacity*2; i++ {
		if cache.Contains(i) {
			t.FailNow()
		}
	}
}

func TestCache_Cap(t *testing.T) {
	capacity := 10
	cache, _ := NewLruCache(capacity)
	if cache.Cap() != capacity {
		t.FailNow()
	}
}

func TestCache_Len(t *testing.T) {
	capacity := 10
	cache, _ := NewLruCache(capacity)
	if cache.Len() != 0 {
		t.FailNow()
	}
	for i := 0; i < capacity*2; i++ {
		cache.Add(i, strconv.Itoa(i))
		if cache.Len() != int(math.Min(float64(i+1), float64(capacity))) {
			t.FailNow()
		}
	}
}

func TestCache_Clear(t *testing.T) {
	capacity := 10
	cache, _ := NewLruCache(capacity)
	for i := 0; i < capacity; i++ {
		cache.Add(i, strconv.Itoa(i))
	}
	cache.Clear()
	if cache.Len() != 0 {
		t.FailNow()
	}
}

func TestCache_SetCap(t *testing.T) {
	capacity := 10
	newCapacity := 20
	cache, _ := NewLruCache(capacity)
	err := cache.SetCap(newCapacity)
	if err != nil {
		t.FailNow()
	}
	if cache.Cap() != newCapacity {
		t.FailNow()
	}
	newCapacity = 5
	err = cache.SetCap(newCapacity)
	if err != nil {
		t.FailNow()
	}
	if cache.Cap() != newCapacity {
		t.FailNow()
	}
	newCapacity = -1
	err = cache.policy.SetCap(newCapacity)
	if err == nil {
		t.FailNow()
	}
	if cache.Cap() == newCapacity {
		t.FailNow()
	}
}

func TestCache_Keys(t *testing.T) {
	capacity := 10
	cache, _ := NewLruCache(capacity)
	for i := 0; i < capacity*2; i++ {
		cache.Add(i, strconv.Itoa(i))
		keys := cache.Keys()
		if !contains(keys, i) {
			t.FailNow()
		}
	}
}

func TestCache_Values(t *testing.T) {
	capacity := 10
	cache, _ := NewLruCache(capacity)
	for i := 0; i < capacity*2; i++ {
		cache.Add(i, strconv.Itoa(i))
		values := cache.Values()
		if !contains(values, strconv.Itoa(i)) {
			t.FailNow()
		}
	}
}

func contains(s []interface{}, e interface{}) bool {
	for _, c := range s {
		if c == e {
			return true
		}
	}
	return false
}
