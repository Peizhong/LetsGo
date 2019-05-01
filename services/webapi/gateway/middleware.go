package gateway

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	log "github.com/peizhong/letsgo/framework/log"
)

type key int

const (
	ErrorK key = iota
	LogK
	TraceK
	ReRouteK
)

type GWContext struct {
	context.Context
	TracingID string
}

func errorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r.WithContext(GWContext{
			r.Context(),
			"",
		}))
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Info(r.RequestURI)
		stratTime := r.Context().Value(ErrorK)
		log.Info("in logging: %v", stratTime)
		wvc := context.WithValue(r.Context(), LogK, "Log")
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r.WithContext(wvc))
	})
}

func tracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		tracingID := uuid.New().String()
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r.WithContext(GWContext{
			r.Context(),
			tracingID,
		}))
	})
}

func reRoutingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		tracingID := r.Context().Value(TraceK)
		log.Info("in rerouting, tracingId: %v", tracingID)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
