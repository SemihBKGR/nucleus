package nucleus

import (
	"math"
	"testing"
)

func TestOfCache(t *testing.T) {
	capacity := 10
	cache, _ := NewLruCache(capacity)
	if cache.Cap() != capacity {
		t.Fatalf("cap() != %d", capacity)
	}
	if cache.Len() != 0 {
		t.Fatalf("len() != %d", cache.Len())
	}
	for i := 0; i < capacity*2; i++ {
		cache.Add(i, i)
		if cache.Len() != int(math.Min(float64(i+1), float64(capacity))) {
			t.Fatalf("len() != %d", i)
		}
	}

	for i := 0; i < capacity*2; i++ {
		ok := cache.Contains(i)
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
			value, ok := cache.Get(i)
			if !ok {
				t.Fatalf("key '%d' does not exist", i)
			}
			if value.(int) != i {
				t.Fatalf("value != %d", i)
			}
		}

	}

}
