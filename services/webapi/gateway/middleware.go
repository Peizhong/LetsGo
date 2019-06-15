package gateway

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	httpclient "letsgo/framework/http"
	log "letsgo/framework/log"
)

type key int

const (
	TraceK key = iota
	ReRouteK
)

type Header struct {
	K, V string
}

type ReRouteInfo struct {
	DestURL     string
	DestHeaders []Header
	RecvData    []byte
}

type GWContext struct {
	context.Context
	TracingID   string
	ReRouteInfo *ReRouteInfo
}

func errorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Error("error", r)
			}
		}()
		// Do stuff here
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Info(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func tracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		tracingID := uuid.New().String()
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		gwContext := GWContext{
			r.Context(),
			tracingID,
			new(ReRouteInfo),
		}
		next.ServeHTTP(w, r.WithContext(gwContext))
	})
}

func reRoutingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		if gwContext, ok := r.Context().(GWContext); ok {
			tracingID := gwContext.TracingID
			log.Info("in rerouting, tracingId: %v", tracingID)
			// reroute
			destURL := ConvertURL(r.RequestURI)
			gwContext.ReRouteInfo.DestURL = destURL
			//
		}
		next.ServeHTTP(w, r)
	})
}

func requestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		if gwContext, ok := r.Context().(GWContext); ok {
			destURL := gwContext.ReRouteInfo.DestURL
			if res, err := httpclient.Do(r.Method, destURL, nil, "", nil); err == nil {
				gwContext.ReRouteInfo.RecvData = res.Body
				gwContext.ReRouteInfo.DestHeaders = make([]Header, len(res.Headers))
				for i, h := range res.Headers {
					gwContext.ReRouteInfo.DestHeaders[i] = Header{
						h.K,
						h.V,
					}
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}

func responseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		if gwContext, ok := r.Context().(GWContext); ok {
			for _, h := range gwContext.ReRouteInfo.DestHeaders {
				w.Header().Set(h.K, h.V)
			}
			w.Write(gwContext.ReRouteInfo.RecvData)
		}
		next.ServeHTTP(w, r)
	})
}
