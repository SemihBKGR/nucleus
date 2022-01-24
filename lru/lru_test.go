package lru

import (
	"math"
	"strconv"
	"testing"
)

func TestNewLru(t *testing.T) {
	lru, err := NewLru(1)
	if lru == nil {
		t.FailNow()
	}
	if err != nil {
		t.FailNow()
	}
	lru, err = NewLru(0)
	if lru != nil {
		t.FailNow()
	}
	if err == nil {
		t.FailNow()
	}
	lru, err = NewLru(-1)
	if lru != nil {
		t.FailNow()
	}
	if err == nil {
		t.FailNow()
	}
}

func TestLru_Add(t *testing.T) {
	capacity := 10
	lru, _ := NewLru(capacity)
	for i := 0; i < capacity*2; i++ {
		lru.Add(i, strconv.Itoa(i))
	}
}

func TestLru_Add2(t *testing.T) {
	capacity := 3
	lru, _ := NewLru(capacity)
	// 1 -> []
	lru.Add(1, nil)
	// [1]
	if !contains(lru.Keys(), 1) {
		t.FailNow()
	}
	// 2 -> [1]
	lru.Add(2, nil)
	// [1,2]
	if !containsAll(lru.Keys(), 1, 2) {
		t.FailNow()
	}
	// 3 -> [1,2]
	lru.Add(3, nil)
	// [1,2,3]
	if !containsAll(lru.Keys(), 1, 2, 3) {
		t.FailNow()
	}
	// 4 -> [1,2,3]
	lru.Add(4, nil)
	// [2,3,4]
	if !containsAll(lru.Keys(), 2, 3, 4) {
		t.FailNow()
	}
	// 2 -> [2,3,4]
	lru.Add(2, nil)
	// [3,4,2]
	if !containsAll(lru.Keys(), 2, 3, 4) {
		t.FailNow()
	}
	// 5 -> [3,4,2]
	lru.Add(5, nil)
	// [4,2,5]
	if !containsAll(lru.Keys(), 2, 4, 5) {
		t.FailNow()
	}
}

func TestLru_Get(t *testing.T) {
	capacity := 10
	lru, _ := NewLru(capacity)
	for i := 0; i < capacity; i++ {
		lru.Add(i, strconv.Itoa(i))
	}
	for i := 0; i < capacity; i++ {
		value, ok := lru.Get(i, true)
		if !ok || value.(string) != strconv.Itoa(i) {
			t.FailNow()
		}
	}
	for i := capacity; i < capacity*2; i++ {
		value, ok := lru.Get(i, true)
		if ok || value != nil {
			t.FailNow()
		}
	}
}

func TestLru_Remove(t *testing.T) {
	capacity := 10
	lru, _ := NewLru(capacity)
	for i := 0; i < capacity; i++ {
		lru.Add(i, strconv.Itoa(i))
	}
	for i := 0; i < capacity; i++ {
		ok := lru.Remove(i)
		if !ok {
			t.FailNow()
		}
	}
	for i := 0; i < capacity; i++ {
		ok := lru.Remove(i)
		if ok {
			t.FailNow()
		}
	}
}

func TestLru_Clear(t *testing.T) {
	capacity := 10
	lru, _ := NewLru(capacity)
	for i := 0; i < capacity; i++ {
		lru.Add(i, strconv.Itoa(i))
	}
	lru.Clear()
	if lru.Len() != 0 {
		t.FailNow()
	}
}

func TestLru_Cap(t *testing.T) {
	capacity := 10
	lru, _ := NewLru(capacity)
	if lru.Cap() != capacity {
		t.FailNow()
	}
}

func TestLru_Len(t *testing.T) {
	capacity := 10
	lru, _ := NewLru(capacity)
	if lru.Len() != 0 {
		t.FailNow()
	}
	for i := 0; i < capacity*2; i++ {
		lru.Add(i, strconv.Itoa(i))
		if lru.Len() != int(math.Min(float64(i+1), float64(capacity))) {
			t.FailNow()
		}
	}
}

func TestLru_SetCap(t *testing.T) {
	capacity := 10
	newCapacity := 20
	lru, _ := NewLru(capacity)
	err := lru.SetCap(newCapacity)
	if err != nil {
		t.FailNow()
	}
	if lru.Cap() != newCapacity {
		t.FailNow()
	}
	newCapacity = 5
	err = lru.SetCap(newCapacity)
	if err != nil {
		t.FailNow()
	}
	if lru.Cap() != newCapacity {
		t.FailNow()
	}
	newCapacity = -1
	err = lru.SetCap(newCapacity)
	if err == nil {
		t.FailNow()
	}
	if lru.Cap() == newCapacity {
		t.FailNow()
	}
}

func TestLru_Keys(t *testing.T) {
	capacity := 10
	lru, _ := NewLru(capacity)
	for i := 0; i < capacity*2; i++ {
		lru.Add(i, strconv.Itoa(i))
		keys := lru.Keys()
		if !contains(keys, i) {
			t.FailNow()
		}
	}
}

func TestLru_Values(t *testing.T) {
	capacity := 10
	lru, _ := NewLru(capacity)
	for i := 0; i < capacity*2; i++ {
		lru.Add(i, strconv.Itoa(i))
		values := lru.Values()
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

func containsAll(s []interface{}, es ...interface{}) bool {
	for _, e := range es {
		if !contains(s, e) {
			return false
		}
	}
	return true
}
