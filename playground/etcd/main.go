package main

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sync"
	"time"

	"github.com/emirpasic/gods/lists"
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/peizhong/letsgo/internal"

	"github.com/coreos/etcd/clientv3"
	etcdnaming "go.etcd.io/etcd/clientv3/naming"

	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer"
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
		cfg := clientv3.Config{
			Endpoints:   []string{"http://193.112.41.28:2379"},
			DialTimeout: 2 * time.Second,
		}
		client, err := clientv3.New(cfg)
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
				if v, ok := value.Interface().(clientv3.WatchResponse); ok {
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

// 将etcd的访问封装成grpc的负载均衡器，grpc调用服务指定域名，etcd负载均衡器解析成实际的ip
func grpccall() {
	cli, cerr := clientv3.NewFromURL("http://localhost:2379")
	if cerr != nil {
		return
	}
	// etcd的域名解析
	r := &etcdnaming.GRPCResolver{Client: cli}
	b := grpc.RoundRobin(r)
	// grpc只用传域名
	conn, gerr := grpc.Dial("my-service", grpc.WithBalancer(b), grpc.WithBlock())
	if gerr != nil {
		return
	}
	defer conn.Close()
	// r.Update(context.TODO(), "my-service", naming.Update{Op: naming.Add, Addr: "1.2.3.4", Metadata: "..."})
}

func main() {
	internal.HostWithContext(app, internal.TimeoutArg{T: time.Second})
}
