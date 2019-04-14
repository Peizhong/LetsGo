package gonet

import (
	"fmt"
	"io"
	"net/http"
	"runtime"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
)

func init() {
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)
}

var once sync.Once

var entryPoint func(*Context) error

var reporter jaeger.Reporter

func configJaeger() {
	sender, err := jaeger.NewUDPTransport("192.168.3.19:5775", 0)
	if err != nil {
		log.Error("Failed to start jaeger udp transport")
	}
	reporter = jaeger.NewRemoteReporter(
		sender,
		jaeger.ReporterOptions.BufferFlushInterval(1*time.Second),
	)
}

func newTracer() (opentracing.Tracer, io.Closer) {
	tracer, closer := jaeger.NewTracer(
		"LetsGo",
		jaeger.NewConstSampler(true),
		reporter)
	return tracer, closer
}

func gatewayHandler(w http.ResponseWriter, req *http.Request) {
	once.Do(func() {
		configJaeger()
		entryPoint = BuildPipeline()
	})
	tracer, closer := newTracer()
	span := tracer.StartSpan("myspan")
	span.SetTag("mytag", "123")
	tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	context := &Context{
		SrcPath:   req.RequestURI[1:],
		Request:   req,
		Responser: w,
		Tracer:    tracer,
	}
	err := entryPoint(context)
	if err != nil {
		fmt.Println("request error")
	}
	span.Finish()
	closer.Close()
}

// RunHTTPServer start gonet http server
func RunHTTPServer(ip string, port int) (err error) {
	address := fmt.Sprintf("%v:%v", ip, port)
	fmt.Println("HTTP server listening at:", address)
	http.HandleFunc("/", gatewayHandler)
	// 替代默认的DefaultServeMux
	//gateway := Gateway{}
	//err = http.ListenAndServe(address, gateway)
	err = http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
	return
}
