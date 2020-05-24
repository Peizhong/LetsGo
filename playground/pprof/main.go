package main

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/peizhong/letsgo/pkg/db"
	"golang.org/x/sync/singleflight"
	"golang.org/x/sys/cpu"

	_ "net/http/pprof"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// _ "go.uber.org/automaxprocs"

type NoPad struct {
	a uint64
	b uint64
	c uint64
}

func (np *NoPad) Increase() {
	atomic.AddUint64(&np.a, 1)
	atomic.AddUint64(&np.b, 1)
	atomic.AddUint64(&np.c, 1)
}

type Pad struct {
	a uint64
	// 数据被多个goroutine访问，如果多个cpu同时访问某个变量，可能cpu把变量和相邻数据都读到缓存中
	// 当cpu1更新变量时，cpu2的变量无效(false sharing)，导致cpu2没用到的变量也更新
	// 将变量补足为64字节，填充一个cacheline
	_p1 [8]uint64
	b   uint64
	_p2 [8]uint64
	c   uint64
	_p3 [8]uint64
}

func (p *Pad) Increase() {
	atomic.AddUint64(&p.a, 1)
	atomic.AddUint64(&p.b, 1)
	atomic.AddUint64(&p.c, 1)
}

type MPad struct {
	_p1 [7]uint64
	a   uint64
	_p2 [7]uint64
	b   uint64
	_p3 [7]uint64
	c   uint64
}

func (mp *MPad) Increase() {
	atomic.AddUint64(&mp.a, 1)
	atomic.AddUint64(&mp.b, 1)
	atomic.AddUint64(&mp.c, 1)
}

type SysPad struct {
	_ cpu.CacheLinePad
	a uint64
	_ cpu.CacheLinePad
	b uint64
	_ cpu.CacheLinePad
	c uint64
}

func (sp *SysPad) Increase() {
	atomic.AddUint64(&sp.a, 1)
	atomic.AddUint64(&sp.b, 1)
	atomic.AddUint64(&sp.c, 1)
}

type database struct {
	database db.ORMHandler
}

func (d *database) Handler(writer http.ResponseWriter, request *http.Request) {
	type money_account struct {
		Id          int    `gorm:"Column:id"`
		AccountName string `gorm:"Column:account_name;type:varchar(200);"`
	}
	var ma money_account
	err := d.database.Get(&ma, db.Query{Key: "id", Op: "=", Value: 2})
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(writer, "%v", ma)
}

// export GODEBUG=gctrace=1

// go tool compile -S main.go

func main() {
	// "net/http/pprof"
	// debug/pprof/profile：访问这个链接会自动进行 CPU profiling，持续 30s，并生成一个文件供下载
	// debug/pprof/heap： Memory Profiling 的路径，访问这个链接会得到一个内存 Profiling 结果的文件，在应用程序进行堆分配时记录堆栈跟踪，用于监视当前和历史内存使用情况，以及检查内存泄漏
	// debug/pprof/block：查看导致阻塞同步的堆栈跟踪
	// debug/pprof/goroutines：运行的 goroutines 列表，以及调用关系

	// https://eddycjy.com/posts/go/tools/2018-09-15-go-tool-pprof/
	// go tool pprof -http :8080 http://localhost:8000/debug/pprof/profile

	http.HandleFunc("/ab", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("OK"))
	})
	flight := &singleflight.Group{}
	http.HandleFunc("/run", func(writer http.ResponseWriter, request *http.Request) {
		var shareCnt int
		for i := 0; i < 100; i++ {
			_, _, shared := flight.Do(fmt.Sprintf("%d", i), func() (i interface{}, e error) {
				pad := SysPad{}
				pad.Increase()
				return pad, nil
			})
			if shared {
				shareCnt++
			}
		}
		fmt.Fprintf(writer, "%d", shareCnt)
	})
	storage := database{db.DBFactory("mysql")}
	http.HandleFunc("/ping", storage.Handler)
	if err := http.ListenAndServe(":8000", nil); err != nil {
		panic(err)
	}
}
