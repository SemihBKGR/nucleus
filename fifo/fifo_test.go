package fifo

import (
	"math"
	"strconv"
	"testing"
)

func TestNewFifo(t *testing.T) {
	fifo, err := NewFifo(1)
	if fifo == nil {
		t.FailNow()
	}
	if err != nil {
		t.FailNow()
	}
	fifo, err = NewFifo(0)
	if fifo != nil {
		t.FailNow()
	}
	if err == nil {
		t.FailNow()
	}
	fifo, err = NewFifo(-1)
	if fifo != nil {
		t.FailNow()
	}
	if err == nil {
		t.FailNow()
	}
}

func TestFifo_Add(t *testing.T) {
	capacity := 10
	fifo, _ := NewFifo(capacity)
	for i := 0; i < capacity*2; i++ {
		fifo.Add(i, strconv.Itoa(i))
	}
}

func TestFifo_Add2(t *testing.T) {
	capacity := 3
	fifo, _ := NewFifo(capacity)
	// 1 -> []
	fifo.Add(1, nil)
	// [1]
	if !contains(fifo.Keys(), 1) {
		t.FailNow()
	}
	// 2 -> [1]
	fifo.Add(2, nil)
	// [1,2]
	if !containsAll(fifo.Keys(), 1, 2) {
		t.FailNow()
	}
	// 3 -> [1,2]
	fifo.Add(3, nil)
	// [1,2,3]
	if !containsAll(fifo.Keys(), 1, 2, 3) {
		t.FailNow()
	}
	// 4 -> [1,2,3]
	fifo.Add(4, nil)
	// [2,3,4]
	if !containsAll(fifo.Keys(), 2, 3, 4) {
		t.FailNow()
	}
	// 2 -> [2,3,4]
	fifo.Add(2, nil)
	// [3,4,2]
	if !containsAll(fifo.Keys(), 2, 3, 4) {
		t.FailNow()
	}
	// 5 -> [3,4,2]
	fifo.Add(5, nil)
	// [[4,2,5]
	if !containsAll(fifo.Keys(), 2, 4, 5) {
		t.FailNow()
	}
}

func TestFifo_Get(t *testing.T) {
	capacity := 10
	fifo, _ := NewFifo(capacity)
	for i := 0; i < capacity; i++ {
		fifo.Add(i, strconv.Itoa(i))
	}
	for i := 0; i < capacity; i++ {
		value, ok := fifo.Get(i, true)
		if !ok || value.(string) != strconv.Itoa(i) {
			t.FailNow()
		}
	}
	for i := capacity; i < capacity*2; i++ {
		value, ok := fifo.Get(i, true)
		if ok || value != nil {
			t.FailNow()
		}
	}
}

func TestFifo_Remove(t *testing.T) {
	capacity := 10
	fifo, _ := NewFifo(capacity)
	for i := 0; i < capacity; i++ {
		fifo.Add(i, strconv.Itoa(i))
	}
	for i := 0; i < capacity; i++ {
		ok := fifo.Remove(i)
		if !ok {
			t.FailNow()
		}
	}
	for i := 0; i < capacity; i++ {
		ok := fifo.Remove(i)
		if ok {
			t.FailNow()
		}
	}
}

func TestFifo_Clear(t *testing.T) {
	capacity := 10
	fifo, _ := NewFifo(capacity)
	for i := 0; i < capacity; i++ {
		fifo.Add(i, strconv.Itoa(i))
	}
	fifo.Clear()
	if fifo.Len() != 0 {
		t.FailNow()
	}
}

func TestFifo_Cap(t *testing.T) {
	capacity := 10
	fifo, _ := NewFifo(capacity)
	if fifo.Cap() != capacity {
		t.FailNow()
	}
}

func TestFifo_Len(t *testing.T) {
	capacity := 10
	fifo, _ := NewFifo(capacity)
	if fifo.Len() != 0 {
		t.FailNow()
	}
	for i := 0; i < capacity*2; i++ {
		fifo.Add(i, strconv.Itoa(i))
		if fifo.Len() != int(math.Min(float64(i+1), float64(capacity))) {
			t.FailNow()
		}
	}
}

func TestFifo_SetCap(t *testing.T) {
	capacity := 10
	newCapacity := 20
	fifo, _ := NewFifo(capacity)
	err := fifo.SetCap(newCapacity)
	if err != nil {
		t.FailNow()
	}
	if fifo.Cap() != newCapacity {
		t.FailNow()
	}
	newCapacity = 5
	err = fifo.SetCap(newCapacity)
	if err != nil {
		t.FailNow()
	}
	if fifo.Cap() != newCapacity {
		t.FailNow()
	}
	newCapacity = -1
	err = fifo.SetCap(newCapacity)
	if err == nil {
		t.FailNow()
	}
	if fifo.Cap() == newCapacity {
		t.FailNow()
	}
}

func TestFifo_Keys(t *testing.T) {
	capacity := 10
	fifo, _ := NewFifo(capacity)
	for i := 0; i < capacity*2; i++ {
		fifo.Add(i, strconv.Itoa(i))
		keys := fifo.Keys()
		if !contains(keys, i) {
			t.FailNow()
		}
	}
}

func TestFifo_Values(t *testing.T) {
	capacity := 10
	fifo, _ := NewFifo(capacity)
	for i := 0; i < capacity*2; i++ {
		fifo.Add(i, strconv.Itoa(i))
		values := fifo.Values()
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
