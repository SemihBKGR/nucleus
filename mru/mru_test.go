package mru

import (
	"math"
	"strconv"
	"testing"
)

func TestNewMru(t *testing.T) {
	mru, err := NewMru(1)
	if mru == nil {
		t.FailNow()
	}
	if err != nil {
		t.FailNow()
	}
	mru, err = NewMru(0)
	if mru != nil {
		t.FailNow()
	}
	if err == nil {
		t.FailNow()
	}
	mru, err = NewMru(-1)
	if mru != nil {
		t.FailNow()
	}
	if err == nil {
		t.FailNow()
	}
}

func TestMru_Add(t *testing.T) {
	capacity := 10
	mru, _ := NewMru(capacity)
	for i := 0; i < capacity*2; i++ {
		mru.Add(i, strconv.Itoa(i))
	}
}

func TestMru_Add2(t *testing.T) {
	capacity := 3
	mru, _ := NewMru(capacity)
	// 1 -> []
	mru.Add(1, nil)
	// [1]
	if !contains(mru.Keys(), 1) {
		t.FailNow()
	}
	// 2 -> [1]
	mru.Add(2, nil)
	// [1,2]
	if !containsAll(mru.Keys(), 1, 2) {
		t.FailNow()
	}
	// 3 -> [1,2]
	mru.Add(3, nil)
	// [1,2,3]
	if !containsAll(mru.Keys(), 1, 2, 3) {
		t.FailNow()
	}
	// 4 -> [1,2,3]
	mru.Add(4, nil)
	// [1,2,4]
	if !containsAll(mru.Keys(), 1, 2, 4) {
		t.FailNow()
	}
	// 2 -> [1,2,4]
	mru.Add(2, nil)
	// [1,4,2]
	if !containsAll(mru.Keys(), 1, 2, 4) {
		t.FailNow()
	}
	// 5 -> [1,4,2]
	mru.Add(5, nil)
	// [1,4,5]
	if !containsAll(mru.Keys(), 1, 4, 5) {
		t.FailNow()
	}
}

func TestMru_Get(t *testing.T) {
	capacity := 10
	mru, _ := NewMru(capacity)
	for i := 0; i < capacity; i++ {
		mru.Add(i, strconv.Itoa(i))
	}
	for i := 0; i < capacity; i++ {
		value, ok := mru.Get(i, true)
		if !ok || value.(string) != strconv.Itoa(i) {
			t.FailNow()
		}
	}
	for i := capacity; i < capacity*2; i++ {
		value, ok := mru.Get(i, true)
		if ok || value != nil {
			t.FailNow()
		}
	}
}

func TestMru_Remove(t *testing.T) {
	capacity := 10
	mru, _ := NewMru(capacity)
	for i := 0; i < capacity; i++ {
		mru.Add(i, strconv.Itoa(i))
	}
	for i := 0; i < capacity; i++ {
		ok := mru.Remove(i)
		if !ok {
			t.FailNow()
		}
	}
	for i := 0; i < capacity; i++ {
		ok := mru.Remove(i)
		if ok {
			t.FailNow()
		}
	}
}

func TestMru_Clear(t *testing.T) {
	capacity := 10
	mru, _ := NewMru(capacity)
	for i := 0; i < capacity; i++ {
		mru.Add(i, strconv.Itoa(i))
	}
	mru.Clear()
	if mru.Len() != 0 {
		t.FailNow()
	}
}

func TestMru_Cap(t *testing.T) {
	capacity := 10
	mru, _ := NewMru(capacity)
	if mru.Cap() != capacity {
		t.FailNow()
	}
}

func TestMru_Len(t *testing.T) {
	capacity := 10
	mru, _ := NewMru(capacity)
	if mru.Len() != 0 {
		t.FailNow()
	}
	for i := 0; i < capacity*2; i++ {
		mru.Add(i, strconv.Itoa(i))
		if mru.Len() != int(math.Min(float64(i+1), float64(capacity))) {
			t.FailNow()
		}
	}
}

func TestMru_SetCap(t *testing.T) {
	capacity := 10
	newCapacity := 20
	mru, _ := NewMru(capacity)
	err := mru.SetCap(newCapacity)
	if err != nil {
		t.FailNow()
	}
	if mru.Cap() != newCapacity {
		t.FailNow()
	}
	newCapacity = 5
	err = mru.SetCap(newCapacity)
	if err != nil {
		t.FailNow()
	}
	if mru.Cap() != newCapacity {
		t.FailNow()
	}
	newCapacity = -1
	err = mru.SetCap(newCapacity)
	if err == nil {
		t.FailNow()
	}
	if mru.Cap() == newCapacity {
		t.FailNow()
	}
}

func TestMru_Keys(t *testing.T) {
	capacity := 10
	mru, _ := NewMru(capacity)
	for i := 0; i < capacity*2; i++ {
		mru.Add(i, strconv.Itoa(i))
		keys := mru.Keys()
		if !contains(keys, i) {
			t.FailNow()
		}
	}
}

func TestMru_Values(t *testing.T) {
	capacity := 10
	mru, _ := NewMru(capacity)
	for i := 0; i < capacity*2; i++ {
		mru.Add(i, strconv.Itoa(i))
		values := mru.Values()
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
