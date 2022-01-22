package nucleus

import (
	"math"
	"strconv"
	"testing"
)

func TestNewLru(t *testing.T) {
	lru, err := newLru(1)
	if lru == nil {
		t.FailNow()
	}
	if err != nil {
		t.FailNow()
	}
	lru, err = newLru(0)
	if lru != nil {
		t.FailNow()
	}
	if err == nil {
		t.FailNow()
	}
	lru, err = newLru(-1)
	if lru != nil {
		t.FailNow()
	}
	if err == nil {
		t.FailNow()
	}
}

func TestLru_add(t *testing.T) {
	capacity := 10
	cache, _ := NewLruCache(capacity)
	for i := 0; i < capacity*2; i++ {
		cache.Add(i, strconv.Itoa(i))
	}
}

func TestLru_get(t *testing.T) {
	capacity := 10
	lru, _ := newLru(capacity)
	for i := 0; i < capacity; i++ {
		lru.add(i, strconv.Itoa(i))
	}
	for i := 0; i < capacity; i++ {
		value, ok := lru.get(i, true)
		if !ok || value.(string) != strconv.Itoa(i) {
			t.FailNow()
		}
	}
	for i := capacity; i < capacity*2; i++ {
		value, ok := lru.get(i, true)
		if ok || value != nil {
			t.FailNow()
		}
	}
}

func TestLru_remove(t *testing.T) {
	capacity := 10
	lru, _ := newLru(capacity)
	for i := 0; i < capacity; i++ {
		lru.add(i, strconv.Itoa(i))
	}
	for i := 0; i < capacity; i++ {
		ok := lru.remove(i)
		if !ok {
			t.FailNow()
		}
	}
	for i := 0; i < capacity; i++ {
		ok := lru.remove(i)
		if ok {
			t.FailNow()
		}
	}
}

func TestLru_contains(t *testing.T) {
	capacity := 10
	lru, _ := newLru(capacity)
	for i := 0; i < capacity; i++ {
		lru.add(i, strconv.Itoa(i))
	}
	for i := 0; i < capacity; i++ {
		if !lru.contains(i) {
			t.FailNow()
		}
	}
	for i := capacity; i < capacity*2; i++ {
		if lru.contains(i) {
			t.FailNow()
		}
	}
}

func TestLru_cap(t *testing.T) {
	capacity := 10
	lru, _ := newLru(capacity)
	if lru.cap() != capacity {
		t.FailNow()
	}
}

func TestLru_len(t *testing.T) {
	capacity := 10
	lru, _ := newLru(capacity)
	if lru.len() != 0 {
		t.FailNow()
	}
	for i := 0; i < capacity*2; i++ {
		lru.add(i, strconv.Itoa(i))
		if lru.len() != int(math.Min(float64(i+1), float64(capacity))) {
			t.FailNow()
		}
	}
}

func TestLru_Purge(t *testing.T) {
	capacity := 10
	lru, _ := newLru(capacity)
	for i := 0; i < capacity; i++ {
		lru.add(i, strconv.Itoa(i))
	}
	lru.purge()
	if lru.len() != 0 {
		t.FailNow()
	}
}

func TestLru_reCap(t *testing.T) {
	capacity := 10
	newCapacity := 20
	lru, _ := newLru(capacity)
	err := lru.reCap(newCapacity)
	if err != nil {
		t.FailNow()
	}
	if lru.cap() != newCapacity {
		t.FailNow()
	}
	newCapacity = 5
	err = lru.reCap(newCapacity)
	if err != nil {
		t.FailNow()
	}
	if lru.cap() != newCapacity {
		t.FailNow()
	}
	newCapacity = -1
	err = lru.reCap(newCapacity)
	if err == nil {
		t.FailNow()
	}
	if lru.cap() == newCapacity {
		t.FailNow()
	}
}

func TestLru_keys(t *testing.T) {
	capacity := 10
	lru, _ := newLru(capacity)
	for i := 0; i < capacity*2; i++ {
		lru.add(i, strconv.Itoa(i))
		keys := lru.keys()
		if !contains(keys, i) {
			t.FailNow()
		}
	}
}

func TestLru_values(t *testing.T) {
	capacity := 10
	cache, _ := newLru(capacity)
	for i := 0; i < capacity*2; i++ {
		cache.add(i, strconv.Itoa(i))
		values := cache.values()
		if !contains(values, strconv.Itoa(i)) {
			t.FailNow()
		}
	}
}
