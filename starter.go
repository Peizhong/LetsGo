package main

import (
	"sync"
	"time"

	"github.com/peizhong/letsgo/framework/log"

	// grpc
	"github.com/peizhong/letsgo/services/grpc/server/helloworld"

	// webapi
	"github.com/peizhong/letsgo/services/webapi/catalogservice"
	"github.com/peizhong/letsgo/services/webapi/gatewayservice"
)

type Starter interface {
	Start()
	Close(wg *sync.WaitGroup)
}

type BatchHttp struct {
	Config interface{}
}

func (b *BatchHttp) Start() {
	go func() {
		catalogservice.Run()
	}()
	go func() {
		gatewayservice.Run()
	}()
}

func (BatchHttp) Close(wg *sync.WaitGroup) {
	go func() {
		time.AfterFunc(2*time.Second, func() {
			wg.Done()
			log.Info("Http down")
		})
	}()
}

type BatchGrpc struct {
	Config interface{}
}

func (b *BatchGrpc) Start() {
	go func() {
		helloworld.Run()
	}()
}

func (BatchGrpc) Close(wg *sync.WaitGroup) {
	go func() {
		time.AfterFunc(1*time.Second, func() {
			wg.Done()
			log.Info("Grpc down")
		})
	}()
}

func loadWorkers() []Starter {
	workers := []Starter{
		new(BatchHttp),
		new(BatchGrpc),
	}
	for _, w := range workers {
		w.Start()
	}
	return workers
}

func stopWorkers(workers ...Starter) {
	var wg sync.WaitGroup
	for _, w := range workers {
		wg.Add(1)
		w.Close(&wg)
	}
	wg.Wait()
}
