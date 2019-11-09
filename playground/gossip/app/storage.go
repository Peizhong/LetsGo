package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coreos/etcd/pkg/fileutil"
	"github.com/google/uuid"
	"github.com/hashicorp/memberlist"
	"go.uber.org/atomic"
	"golang.org/x/sync/singleflight"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

type dataType struct {
	Id    []byte
	Value interface{}
	Time  int64
}

type dataStore map[string]dataType

// kv变更的记录，附载在broadcast中传输
type update struct {
	Action string // add, del
	Data   dataStore
}

type storage struct {
	flight     *singleflight.Group
	data       dataStore
	count      atomic.Int32
	lock       sync.RWMutex
	broadcasts *memberlist.TransmitLimitedQueue
	ch         chan *update

	Name      string
	LocalNode func() *memberlist.Node
	Members   func() []string
}

func NewStorage(name string, port int, join string) *storage {
	s := &storage{
		flight: &singleflight.Group{},
		data:   make(dataStore),
		ch:     make(chan *update, 100),
	}
	if name == "" {
		// auto name
		name = func() string {
			host, _ := os.Hostname()
			name := fmt.Sprintf("%s@%s", uuid.New().String()[:7], host)
			return name
		}()
	}
	s.RegisterMember(name, port, join)
	return s
}

func (s *storage) RegisterMember(name string, port int, join string) {
	s.Name = name
	config := memberlist.DefaultLocalConfig()
	config.Name = s.Name
	config.BindPort = port
	config.Events = &event{}
	config.Delegate = &delegate{storage: s}
	member, err := memberlist.Create(config)
	if err != nil {
		log.Panic(err)
	}
	s.broadcasts = &memberlist.TransmitLimitedQueue{
		NumNodes: func() int {
			return member.NumMembers()
		},
		RetransmitMult: 3,
	}
	s.LocalNode = func() *memberlist.Node {
		return member.LocalNode()
	}
	s.Members = func() []string {
		var resp []string
		for _, m := range member.Members() {
			resp = append(resp, fmt.Sprintf("%s@%s:%d", m.Name, m.Addr, m.Port))
		}
		return resp
	}
	if join != "" {
		_, err := member.Join([]string{join})
		if err != nil {
			log.Panic(err)
		}
	}
	// delegate 收到的消息发送到ch
	go s.updateHandler()
}

func (s *storage) Get(key string) (interface{}, bool) {
	// suppress duplicate requests with singleflight)
	v, err, _ := s.flight.Do(key, func() (i interface{}, e error) {
		s.lock.Lock()
		defer s.lock.Unlock()
		if v, ok := s.data[key]; ok {
			return v.Value, nil
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
	if _, ok := s.data[key]; !ok {
		s.count.Inc()
	}
	data := dataType{
		Id:    []byte(uuid.New().String()),
		Value: value,
		Time:  time.Now().Unix(),
	}
	s.data[key] = data
	s.lock.Unlock()
	if b, err := json.Marshal([]*update{
		{
			Action: "add",
			Data: dataStore{
				key: data,
			},
		},
	}); err == nil {
		s.broadcasts.QueueBroadcast(&broadcast{
			msg:    append([]byte("d"), b...),
			notify: nil,
		})
	}
}

func (s *storage) Delete(key string) {
	s.flight.Do(key, func() (i interface{}, e error) {
		s.lock.Lock()
		if _, ok := s.data[key]; ok {
			delete(s.data, key)
			s.count.Dec()
		}
		s.lock.Unlock()
		if b, err := json.Marshal([]*update{
			{
				Action: "del",
				Data: dataStore{
					key: dataType{},
				},
			},
		}); err == nil {
			s.broadcasts.QueueBroadcast(&broadcast{
				msg:    append([]byte("d"), b...),
				notify: nil,
			})
		}
		return nil, nil
	})
}

func (s *storage) Gets() []string {
	s.lock.RLock()
	var r []string
	for k, v := range s.data {
		r = append(r, fmt.Sprintf("%v_%v", k, v.Value))
	}
	s.lock.RUnlock()
	return r
}

// Save: 退出时保存
func (s *storage) Save() (err error) {
	s.lock.Lock()
	l := len(s.data)
	data, err := json.Marshal(s.data)
	s.lock.Unlock()
	if err != nil {
		return err
	}
	log.Println(fmt.Sprintf("save %d data", l))
	err = ioutil.WriteFile(s.LocalNode().Name, data, os.ModePerm)
	return
}

// Load: 启动时加载
func (s *storage) Load() (err error) {
	fileName := s.LocalNode().Name
	if !fileutil.Exist(fileName) {
		return
	}
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	if len(data) > 0 {
		s.lock.Lock()
		err = json.Unmarshal(data, &s.data)
		// todo 本地数据与网络数据合并
		if err == nil {
			log.Println(fmt.Sprintf("recover %d data", len(s.data)))
		}
		s.lock.Unlock()
	}
	return
}

func (s *storage) UnpackNotify(b []byte) {
	var updates []*update
	if err := json.Unmarshal(b, &updates); err != nil {
		return
	}
	for _, u := range updates {
		s.ch <- u
	}
}

func (s *storage) updateHandler() {
	for u := range s.ch {
		s.lock.Lock()
		for k, v := range u.Data {
			if lv, ok := s.data[k]; ok {
				// 比较时间
				if lv.Time > v.Time {
					continue
				}
			}
			switch u.Action {
			case "add":
				s.data[k] = v
			case "del":
				delete(s.data, k)
			}
		}
		s.lock.Unlock()
	}
}

func (s *storage) Benchmark(n int) {
	testData := make([]string, n)
	for i := 0; i < n; i++ {
		testData[i] = uuid.New().String()
	}
	start := time.Now()
	for _, k := range testData {
		s.Set(k, "")
		s.Delete(k)
	}
	log.Println("benchmark:", time.Since(start))
}
