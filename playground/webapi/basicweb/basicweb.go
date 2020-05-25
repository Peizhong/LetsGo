package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"

	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/opentracing/opentracing-go"
	tracinglog "github.com/opentracing/opentracing-go/log"
	"github.com/peizhong/letsgo/pkg/config"
	"github.com/peizhong/letsgo/pkg/db"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/uber/jaeger-client-go"

	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	jprometheus "github.com/uber/jaeger-lib/metrics/prometheus"

	"github.com/peizhong/letsgo/internal"
	"github.com/peizhong/letsgo/playground/webapi/basicweb/models"
)

var (
	srv *http.Server
)

type basicHandler struct {
	responseText string
	counter      func()
}

func (th *basicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	span := opentracing.StartSpan("Hi")
	defer span.Finish()
	// span is too large
	// big := make([]byte, 1024*1024*1024)
	span.LogFields(tracinglog.String("hello", "you"))
	if th.counter != nil {
		th.counter()
	}
	fmt.Fprintf(w, th.responseText)
}

type queryHandler struct {
}

func (th *queryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	span := opentracing.StartSpan("query")
	defer span.Finish()
	// span is too large
	// big := make([]byte, 1024*1024*1024)
	span.LogFields(tracinglog.String("hello", "you"))
	orm := db.DBFactory("mysql")
	account := &models.MoneyAccount{}
	orm.Get(&account)
	fmt.Fprintf(w, account.AccountName)
}

func Start() {
	requestMetrics := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "basic_web",
			Subsystem: "request",
			Name:      "request_total",
			Help:      "Total number of jobs processed by the workers",
		},
		// We will want to monitor the worker ID that processed the
		// job, and the type of job that was processed
		[]string{"route"},
	)
	mux := http.NewServeMux()
	mux.Handle("/", &basicHandler{responseText: "hello world", counter: func() { requestMetrics.WithLabelValues("/").Inc() }})
	mux.Handle("/data", &queryHandler{})
	mux.Handle("/help", &basicHandler{responseText: "help me"})
	mux.Handle("/metrics", promhttp.Handler())
	prometheus.MustRegister(requestMetrics)
	srv = &http.Server{
		Addr:    net.JoinHostPort("", strconv.Itoa(config.HTTPApp1Port)),
		Handler: mux,
	}
	log.Println("basic web service is on", config.HTTPApp1Port)
	if err := srv.ListenAndServe(); err != nil {
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
			LogSpans:           true,
			LocalAgentHostPort: "localhost:6831",
		},
	}
	jLogger := jaegerlog.StdLogger
	jMetricsFactory := jprometheus.New()
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

	internal.PProf(Start, Stop)
}
