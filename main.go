package main

import (
	"letsgo/framework/config"
	"letsgo/framework/log"
	"letsgo/services/webapi/basicweb"
	"letsgo/services/webapi/eventloop"
	"letsgo/services/webapi/gateway"
	"os"
	"os/signal"
	"reflect"
	"sync"
)

type Starter interface {
	Start()
	Stop()
}

var services []Starter

func startServices() {
	services = []Starter{
		&gateway.GatewayService{},
		&basicweb.BasicWeb{},
		&eventloop.EventLoop{},
	}
	for _, s := range services {
		srv := s
		go func() {
			v := reflect.TypeOf(srv).Elem()
			n := v.Name()
			log.Infof("starting %s", n)
			srv.Start()
		}()
	}
	log.Info("all service are running")
}

func stopServices() {
	var wg sync.WaitGroup
	wg.Add(len(services))
	for _, s := range services {
		srv := s
		go func() {
			v := reflect.TypeOf(srv).Elem()
			n := v.Name()
			log.Infof("starting %s", n)
			srv.Stop()
			wg.Done()
		}()
	}
	wg.Wait()
}

func main() {
	log.Info("Let's Go start at", config.RunMode)
	startServices()
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Info("receive close signal")
	stopServices()
	log.Info("Let's Go exit")
}
