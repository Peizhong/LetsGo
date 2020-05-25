package internal

import (
	"context"
	"log"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"time"
)

// for type assert
type Arg interface {
}

type TimeoutArg struct {
	Arg
	T time.Duration
}

type TempArg struct {
	Arg
	Value int
}

// 强制退出/自动运行结束
func Host(start, stop func()) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	_, file, _, _ := runtime.Caller(1)
	log.Println("start ", file)
	// start是个循环
	go start()
	<-c
	log.Println("closing...")
	if stop != nil {
		stop()
	}
	log.Println("bye")
}

func getTimeoutArg(args ...Arg) TimeoutArg {
	for i := 0; i < len(args); i++ {
		if a, ok := args[i].(TimeoutArg); ok {
			return a
		}
	}
	return TimeoutArg{}
}

func HostWithContext(app func(context.Context, chan struct{}), args ...Arg) {
	ctx, cancel := context.WithCancel(context.Background())
	exit := make(chan struct{}, 1)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	log.Println("app start")
	go app(ctx, exit)
	timeoutCh := make(chan struct{})
	timeout := getTimeoutArg(args...)
	// 可能正常退出，或强制退出
loop:
	select {
	// 强制退出
	case <-c:
		log.Println("force closing...")
		// 触发ctx的cancel
		cancel()
		// 开启超时检测
		if timeout.T > 0 {
			go func() {
				<-time.After(timeout.T)
				timeoutCh <- struct{}{}
			}()
		}
		// 下个循环，等待app结束
		goto loop
	case <-timeoutCh:
		log.Println("waiting timeout, whatever")
	// 正常结束
	case <-exit:
		log.Println("app exit")
	}
}

func PProf(start, stop func()) {
	const cpuprofile = "cpu.pprof"
	cpufile, err := os.Create(cpuprofile)
	if err != nil {
		panic(err)
	}
	defer cpufile.Close()

	pprof.StartCPUProfile(cpufile)
	defer func() {
		pprof.StopCPUProfile()
		log.Println("go tool pprof -http=:8080 -no_browser [program]", cpuprofile)
	}()
	Host(start, stop)
}

func Trace(start, stop func()) {
	const traceout = "trace.out"
	tracefile, err := os.Create(traceout)
	if err != nil {
		panic(err)
	}
	defer tracefile.Close()

	trace.Start(tracefile)
	defer func() {
		trace.Stop()
		log.Println(" go tool trace", traceout)
	}()

	Host(start, stop)
}
