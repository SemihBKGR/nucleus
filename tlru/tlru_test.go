package tlru

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestNewTlru(t *testing.T) {
	tlru, err := NewTlru(10, 1*time.Second)
	if err != nil {
		t.FailNow()
	}
	lock := sync.RWMutex{}
	tlru.StartEvictionDaemon(&lock)
	lock.Lock()
	tlru.Add(1, nil)
	tlru.Add(2, nil)
	tlru.Add(3, nil)
	tlru.Add(4, nil)
	lock.Unlock()
	time.Sleep(2 * time.Second)
	lock.Lock()
	tlru.Add(4, nil)
	tlru.Add(5, nil)
	tlru.Add(6, nil)
	lock.Unlock()
	fmt.Println(tlru.Keys())
}
