package internal

import (
	"context"
	"log"
	"sync"

	discovery "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	server "github.com/envoyproxy/go-control-plane/pkg/server/v2"
)

type Callbacks struct {
	Signal         chan struct{}
	Debug          bool
	Fetches        int
	Requests       int
	DeltaRequests  int
	DeltaResponses int
	mu             sync.Mutex
}

var _ server.Callbacks = new(Callbacks)

func (cb *Callbacks) OnStreamOpen(_ context.Context, id int64, typ string) error {
	if cb.Debug {
		log.Printf("stream %d open for %s\n", id, typ)
	}
	return nil
}

func (cb *Callbacks) OnStreamRequest(id int64, req *discovery.DiscoveryRequest) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.Requests++
	if cb.Debug {
		log.Println("on stream request", id, req.TypeUrl)
	}
	if cb.Signal != nil {
		close(cb.Signal)
		cb.Signal = nil
	}
	return nil
}

func (cb *Callbacks) OnStreamResponse(id int64, req *discovery.DiscoveryRequest, resp *discovery.DiscoveryResponse) {
	if cb.Debug {
		log.Println("on stream response", id, resp.TypeUrl, len(resp.Resources))
	}
}

func (cb *Callbacks) OnStreamClosed(id int64) {
	if cb.Debug {
		log.Printf("stream %d closed\n", id)
	}
}

func (cb *Callbacks) OnFetchRequest(_ context.Context, req *discovery.DiscoveryRequest) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.Fetches++
	if cb.Debug {
		log.Println("on fetch request", req.ResourceNames)
	}
	if cb.Signal != nil {
		close(cb.Signal)
		cb.Signal = nil
	}
	return nil
}

func (cb *Callbacks) OnFetchResponse(req *discovery.DiscoveryRequest, resp *discovery.DiscoveryResponse) {
	if cb.Debug {
		log.Println("on fetch response", resp.TypeUrl)
	}
}
