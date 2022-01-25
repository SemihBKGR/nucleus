package tlru

import (
	"math"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestTlru_Add(t *testing.T) {
	capacity := 10
	duration := time.Duration(0)
	tlru, _ := NewTlru(capacity, duration)
	for i := 0; i < capacity*2; i++ {
		tlru.Add(i, strconv.Itoa(i))
	}
}

func TestTlru_Add2(t *testing.T) {
	capacity := 10
	duration := time.Duration(0)
	tlru, _ := NewTlru(capacity, duration)
	// 1 -> []
	tlru.Add(1, nil)
	// [1]
	if !contains(tlru.Keys(), 1) {
		t.FailNow()
	}
	// 2 -> [1]
	tlru.Add(2, nil)
	// [1,2]
	if !containsAll(tlru.Keys(), 1, 2) {
		t.FailNow()
	}
	// 3 -> [1,2]
	tlru.Add(3, nil)
	// [1,2,3]
	if !containsAll(tlru.Keys(), 1, 2, 3) {
		t.FailNow()
	}
	// 4 -> [1,2,3]
	tlru.Add(4, nil)
	// [2,3,4]
	if !containsAll(tlru.Keys(), 2, 3, 4) {
		t.FailNow()
	}
	// 2 -> [2,3,4]
	tlru.Add(2, nil)
	// [3,4,2]
	if !containsAll(tlru.Keys(), 2, 3, 4) {
		t.FailNow()
	}
	// 5 -> [3,4,2]
	tlru.Add(5, nil)
	// [4,2,5]
	if !containsAll(tlru.Keys(), 2, 4, 5) {
		t.FailNow()
	}
}

func TestTlru_StartDaemon(t *testing.T) {
	capacity := 10
	duration := time.Duration(0)
	tlru, _ := NewTlru(capacity, duration)
	lock := sync.RWMutex{}
	ok := tlru.StartDaemon(&lock)
	if ok {
		t.FailNow()
	}
	duration = 10 * time.Second
	tlru, _ = NewTlru(capacity, duration)
	ok = tlru.StartDaemon(&lock)
	if !ok {
		t.FailNow()
	}
	ok = tlru.StartDaemon(&lock)
	if ok {
		t.FailNow()
	}
}

func TestTlru_StartDaemon2(t *testing.T) {
	capacity := 5
	duration := 10 * time.Millisecond
	tlru, _ := NewTlru(capacity, duration)
	lock := sync.RWMutex{}
	tlru.StartDaemon(&lock)
	lock.Lock()
	tlru.Add(1, nil)
	tlru.Add(2, nil)
	tlru.Add(3, nil)
	lock.Unlock()
	time.Sleep(20 * time.Millisecond)
	lock.Lock()
	tlru.Add(4, nil)
	tlru.Add(5, nil)
	tlru.Add(6, nil)
	if tlru.Len() != 3 {
		t.FailNow()
	}
	lock.Unlock()
}

func TestTlru_Get(t *testing.T) {
	capacity := 10
	duration := time.Duration(0)
	tlru, _ := NewTlru(capacity, duration)
	for i := 0; i < capacity; i++ {
		tlru.Add(i, strconv.Itoa(i))
	}
	for i := 0; i < capacity; i++ {
		value, ok := tlru.Get(i, true)
		if !ok || value.(string) != strconv.Itoa(i) {
			t.FailNow()
		}
	}
	for i := capacity; i < capacity*2; i++ {
		value, ok := tlru.Get(i, true)
		if ok || value != nil {
			t.FailNow()
		}
	}
}

func TestTlru_Remove(t *testing.T) {
	capacity := 10
	duration := time.Duration(0)
	tlru, _ := NewTlru(capacity, duration)
	for i := 0; i < capacity; i++ {
		tlru.Add(i, strconv.Itoa(i))
	}
	for i := 0; i < capacity; i++ {
		ok := tlru.Remove(i)
		if !ok {
			t.FailNow()
		}
	}
	for i := 0; i < capacity; i++ {
		ok := tlru.Remove(i)
		if ok {
			t.FailNow()
		}
	}
}

func TestTlru_Clear(t *testing.T) {
	capacity := 10
	duration := time.Duration(0)
	tlru, _ := NewTlru(capacity, duration)
	for i := 0; i < capacity; i++ {
		tlru.Add(i, strconv.Itoa(i))
	}
	tlru.Clear()
	if tlru.Len() != 0 {
		t.FailNow()
	}
}

func TestTlru_Cap(t *testing.T) {
	capacity := 10
	duration := time.Duration(0)
	tlru, _ := NewTlru(capacity, duration)
	if tlru.Cap() != capacity {
		t.FailNow()
	}
}

func TestTlru_Len(t *testing.T) {
	capacity := 10
	duration := time.Duration(0)
	tlru, _ := NewTlru(capacity, duration)
	if tlru.Len() != 0 {
		t.FailNow()
	}
	for i := 0; i < capacity*2; i++ {
		tlru.Add(i, strconv.Itoa(i))
		if tlru.Len() != int(math.Min(float64(i+1), float64(capacity))) {
			t.FailNow()
		}
	}
}

func TestTlru_SetCap(t *testing.T) {
	capacity := 10
	duration := time.Duration(0)
	newCapacity := 20
	tlru, _ := NewTlru(capacity, duration)
	err := tlru.SetCap(newCapacity)
	if err != nil {
		t.FailNow()
	}
	if tlru.Cap() != newCapacity {
		t.FailNow()
	}
	newCapacity = 5
	err = tlru.SetCap(newCapacity)
	if err != nil {
		t.FailNow()
	}
	if tlru.Cap() != newCapacity {
		t.FailNow()
	}
	newCapacity = -1
	err = tlru.SetCap(newCapacity)
	if err == nil {
		t.FailNow()
	}
	if tlru.Cap() == newCapacity {
		t.FailNow()
	}
}

func TestTlru_Keys(t *testing.T) {
	capacity := 10
	duration := time.Duration(0)
	tlru, _ := NewTlru(capacity, duration)
	for i := 0; i < capacity*2; i++ {
		tlru.Add(i, strconv.Itoa(i))
		keys := tlru.Keys()
		if !contains(keys, i) {
			t.FailNow()
		}
	}
}

func TestTlru_Values(t *testing.T) {
	capacity := 10
	duration := time.Duration(0)
	tlru, _ := NewTlru(capacity, duration)
	for i := 0; i < capacity*2; i++ {
		tlru.Add(i, strconv.Itoa(i))
		values := tlru.Values()
		if !contains(values, strconv.Itoa(i)) {
			t.FailNow()
		}
	}
}

func TestTlru_DaemonStarted(t *testing.T) {
	capacity := 10
	duration := 10 * time.Second
	tlru, _ := NewTlru(capacity, duration)
	if tlru.DaemonStarted() {
		t.FailNow()
	}
	tlru.StartDaemon(&sync.RWMutex{})
	if !tlru.DaemonStarted() {
		t.FailNow()
	}
}

func TestTlru_ExpirationDuration(t *testing.T) {
	capacity := 10
	duration := time.Duration(0)
	tlru, _ := NewTlru(capacity, duration)
	if tlru.ExpirationDuration() != duration {
		t.FailNow()
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
