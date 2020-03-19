package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/opentracing/opentracing-go"
	tracinglog "github.com/opentracing/opentracing-go/log"
	"github.com/peizhong/letsgo/pkg/config"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-lib/metrics"

	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"

	"github.com/peizhong/letsgo/internal"
)

var (
	srv *http.Server
)

type basicHandler struct {
	responseText string
}

func (th *basicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	span := opentracing.StartSpan("Hi")
	defer span.Finish()
	// span is too large
	big := make([]byte, 1024*1024*1024)
	span.LogFields(tracinglog.String("wa", string(big)))
	fmt.Fprintf(w, th.responseText)
}

func Start() {
	mux := http.NewServeMux()
	mux.Handle("/", &basicHandler{responseText: "hello world"})
	mux.Handle("/help", &basicHandler{responseText: "help me"})
	srv = &http.Server{
		Addr:    net.JoinHostPort("", strconv.Itoa(config.HTTPApp1Port)),
		Handler: mux,
	}
	log.Println("basic web service is on", config.HTTPApp1Port)
	if err := srv.ListenAndServeTLS(config.CertCrt, config.CertKey); err != nil {
		log.Println(err.Error())
	}
}

func Stop() {
	if err := srv.Shutdown(context.Background()); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}
}

func configTracing() func() {
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
			LocalAgentHostPort:"",
		},
	}
	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory
	closer, err := cfg.InitGlobalTracer(
		"serviceName",
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)
	if err != nil {
		log.Printf("Could not initialize jaeger tracer: %s", err.Error())
		return nil
	}
	return func() {
		closer.Close()
	}
}

func main() {
	closer := configTracing()
	defer closer()

	internal.Host(Start, Stop)
}
