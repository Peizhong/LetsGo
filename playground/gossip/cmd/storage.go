package cmd

import (
	"errors"
	"fmt"
	"golang.org/x/sync/singleflight"
	"sync"
)

type storage struct {
	flight *singleflight.Group
	data   map[string]interface{}
	lock   sync.RWMutex
}

func NewStorage() *storage {
	s := &storage{flight: &singleflight.Group{}}
	return s
}

func (s *storage) Get(key string) (interface{}, bool) {
	v, err, _ := s.flight.Do(key, func() (i interface{}, e error) {
		if v, ok := s.data[key]; ok {
			return v, nil
		}
		return nil, errors.New(key)
	})
	if err != nil {
		return nil, false
	}
	return v, true
}

func (s *storage) Set(key string, value interface{}) {
	s.lock.Lock()
	s.data[key] = value
	s.lock.Unlock()
}

func (s *storage) Gets() []string {
	s.lock.RLock()
	r := []string{}
	for k, v := range s.data {
		r = append(r, fmt.Sprintf("%s_%v", k, v))
	}
	s.lock.RUnlock()
	return r
}
