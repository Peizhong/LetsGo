package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/memberlist"
	"go.uber.org/atomic"
	"golang.org/x/sync/singleflight"
)

type Storage struct {
	flight     *singleflight.Group
	data       dataStore
	count      atomic.Int32
	lock       sync.RWMutex
	broadcasts *memberlist.TransmitLimitedQueue
	ch         chan *update

	Name      string
	LocalNode func() *memberlist.Node
	Members   func() ([]*MemberInfo, error)
}

func NewStorage(name string, port int, join string) *Storage {
	s := &Storage{
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
	// gossip收到的消息发送到ch处理
	go s.updateHandler()
	// 定时保存数据到硬盘
	go func() {
		for {
			select {
			case <-time.After(time.Minute):
				s.Save("")
			}
		}
	}()
	return s
}

func (s *Storage) RegisterMember(name string, port int, join string) {
	s.Name = name
	config := memberlist.DefaultLocalConfig()
	config.Name = s.Name
	config.BindPort = port
	config.Events = &event{}
	config.Delegate = &delegate{storage: s}
	config.Alive = &alive{}
	// config.EnableCompression = true
	// config.Logger = log.New(os.Stderr, "", log.LstdFlags)
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
	s.Members = func() ([]*MemberInfo, error) {
		var resp []*MemberInfo
		for _, m := range member.Members() {
			resp = append(resp, &MemberInfo{
				Name: m.Name,
				Addr: fmt.Sprintf("%s:%d", m.Addr, m.Port),
			})
		}
		return resp, nil
	}
	if join != "" {
		_, err := member.Join([]string{join})
		if err != nil {
			log.Panic(err)
		}
	}
}

func (s *Storage) Get(key string) (interface{}, error) {
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
		return nil, err
	}
	return v, nil
}

func (s *Storage) Set(key string, value interface{}) error {
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
			Action: ActionAdd,
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
	return nil
}

func (s *Storage) Delete(key string) error {
	s.flight.Do(key, func() (i interface{}, e error) {
		s.lock.Lock()
		if _, ok := s.data[key]; ok {
			delete(s.data, key)
			s.count.Dec()
		}
		s.lock.Unlock()
		if b, err := json.Marshal([]*update{
			{
				Action: ActionDel,
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
	return nil
}

// Save: 退出时保存
func (s *Storage) Save(path string) (err error) {
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
func (s *Storage) Load(path string) (err error) {
	fileName := s.LocalNode().Name
	data, err := ioutil.ReadFile(fileName)
	if err == nil && len(data) > 0 {
		s.lock.Lock()
		err = json.Unmarshal(data, &s.data)
		// todo: 本地数据与网络数据合并
		if err == nil {
			log.Println(fmt.Sprintf("recover %d data", len(s.data)))
		}
		s.lock.Unlock()
	}
	// 如果是文件不存在，忽略
	if _, ok := err.(*os.PathError); ok {
		return nil
	}
	return err
}

// UnpackNotify: 二次解析，发送到ch
func (s *Storage) UnpackNotify(b []byte) {
	var updates []*update
	if err := json.Unmarshal(b, &updates); err != nil {
		return
	}
	for _, u := range updates {
		s.ch <- u
	}
}

// updateHandler: 处理管道
func (s *Storage) updateHandler() {
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
			case ActionAdd:
				s.data[k] = v
			case ActionDel:
				delete(s.data, k)
			}
		}
		s.lock.Unlock()
	}
}

func (s *Storage) Info() ([]*MemberInfo, error) {
	return s.Members()
}

func (s *Storage) Benchmark(n int) error {
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
	return nil
}
