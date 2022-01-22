package nucleus

import (
	"math"
	"testing"
)

func Test(t *testing.T) {
	capacity := 10
	lru, _ := newLruCache(capacity)
	if lru.cap() != capacity {
		t.Fatalf("cap() != %d", capacity)
	}
	if lru.len() != 0 {
		t.Fatalf("len() != %d", lru.len())
	}
	for i := 0; i < capacity*2; i++ {
		lru.add(i, i)
		if lru.len() != int(math.Min(float64(i+1), float64(capacity))) {
			t.Fatalf("len() != %d", i)
		}
	}

	for i := 0; i < capacity*2; i++ {
		ok := lru.contains(i)
		if i/capacity < 1 {
			if ok {
				t.Fatalf("key '%d' must not be exists", i)
			}
		} else {
			if !ok {
				t.Fatalf("key '%d' must be exists", i)
			}
		}

		for i := capacity; i < capacity*2; i++ {
			value, ok := lru.get(i, false)
			if !ok {
				t.Fatalf("key '%d' does not exist", i)
			}
			if value.(int) != i {
				t.Fatalf("value != %d", i)
			}
		}

	}

}
