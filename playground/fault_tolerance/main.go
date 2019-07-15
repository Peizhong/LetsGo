package main

import (
	"errors"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/sirupsen/logrus"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

var actualJob int64

func job(threadId int) error {
	//log.Println(threadId, "receive")
	exp := hystrix.Go("my_command", func() error {
		log.Println(threadId, "do")
		// talk to other services
		atomic.AddInt64(&actualJob, 1)
		time.Sleep(1 * time.Second)
		return errors.New("not ok")
	}, func(err error) error {
		// log.Println(err)
		switch err {
		case hystrix.ErrCircuitOpen:
			return err
		case hystrix.ErrMaxConcurrency:
			return err
		case hystrix.ErrTimeout:
			return err
		}
		return errors.New("callback")
	})
	select {
	case _ = <-exp:
		log.Println(threadId)
	}
	return nil
}

func producer(total, size int) <-chan int {
	ch := make(chan int, size)
	go func() {
		for i := 0; i < total; i++ {
			ch <- i
		}
		close(ch)
	}()
	return ch
}

func main() {
	hystrix.ConfigureCommand("my_command", hystrix.CommandConfig{
		// Timeout:               1000,
		MaxConcurrentRequests: 1000,
		ErrorPercentThreshold: 10,
		SleepWindow:           1000,
	})
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	hystrix.SetLogger(logger)
	settings := hystrix.GetCircuitSettings()
	log.Println(settings)
	breaker, ok, err := hystrix.GetCircuit("my_command")
	todo := producer(1000, 10)
	var wait sync.WaitGroup
	log.Println(runtime.GOMAXPROCS(0))
	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		wait.Add(1)
		//v := i
		go func() {
			// log.Println("run thread", v)
			for n := range todo {
				if breaker.IsOpen() {
					breaker.ReportEvent()
					log.Println("breaker on, sleep+")
					time.Sleep(1 * time.Second)
				}
				// log.Println("thread", v, "get one job")
				job(n)
			}
			wait.Done()
		}()
	}
	wait.Wait()
	log.Println(breaker.IsOpen())
	log.Println(ok, err)
	log.Println("actual job", actualJob)
}
