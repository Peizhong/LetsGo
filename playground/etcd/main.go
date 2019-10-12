package main

import (
	"context"
	"fmt"
	"github.com/emirpasic/gods/lists"
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/peizhong/letsgo/internal"
	etcd "go.etcd.io/etcd/clientv3"
	"log"
	"reflect"
	"sync"
	"time"
)

type watchList struct {
	mu      sync.Mutex
	toWatch lists.List
}

func (w *watchList) addWatch(key string) *watchList {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.toWatch == nil {
		w.toWatch = arraylist.New()
	}
	if !w.toWatch.Contains(key) {
		w.toWatch.Add(key)
	}
	return w
}

func (w *watchList) watching(ctx context.Context) {
	go func() {
		cfg := etcd.Config{
			Endpoints:   []string{"http://193.112.41.28:2379"},
			DialTimeout: 2 * time.Second,
		}
		client, err := etcd.New(cfg)
		if err != nil {
			log.Println(err)
			return
		}
		defer client.Close()

		cases := make([]reflect.SelectCase, w.toWatch.Size())
		for i, v := range w.toWatch.Values() {
			key := v.(string)
			if g, err := client.Get(ctx, key); err == nil {
				if g.Count == 0 {
					if _, err := client.Put(ctx, key, ""); err != nil {
						fmt.Println(err)
					}
				}
			}
			ch := client.Watch(ctx, key)
			cases[i] = reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(ch),
			}
		}
		cases = append(cases, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ctx.Done()),
		})
		for {
			if _, value, ok := reflect.Select(cases); ok {
				if v, ok := value.Interface().(etcd.WatchResponse); ok {
					for _, e := range v.Events {
						fmt.Println(fmt.Sprintf("op: %v, k: %s, v: %s, rv: %d", e.Type, string(e.Kv.Key), string(e.Kv.Value), e.Kv.ModRevision))
					}
				}
			}
		}
	}()
}

func app(ctx context.Context, exit chan struct{}) {
	wl := watchList{}
	wl.addWatch("ddd").
		addWatch("jiji").
		addWatch("pipi").
		watching(ctx)

	select {
	case <-ctx.Done():
	}
	// todo: 没有等待watching结束

}

func main() {
	internal.HostWithContext(app, internal.TimeoutArg{T: time.Second})
}
