package file

import (
	"sync"
	"testing"
	"time"
)

func Test_changed(t *testing.T) {
	f := Load("D://a.txt")
	wg := new(sync.WaitGroup)
	wg.Add(2)
	f.RegisterFileChanged(func() {
		wg.Done()
	})
	ch := make(chan struct{}, 0)
	go func() {
		wg.Wait()
		ch <- struct{}{}
	}()
	select {
	case <-ch:
		t.Log("Ok")
	case <-time.After(time.Second * 10):
		t.Error("time out")
	}
}

func Test_Select(t *testing.T) {
	var ch1 chan int
	ch1 = make(chan int, 1)
	for i := 0; i < 2; i++ {
		select {
		case ch1 <- 1:
			t.Log("write something")
		case <-ch1:
			t.Log("read something")
		default:
			t.Log("full")
		}
	}
}
