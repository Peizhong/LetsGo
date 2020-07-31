package proxy

import (
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/peizhong/letsgo/pkg/log"
)

// selector select upstream pod
type Selector interface {
	ServiceName() string
	// SelectEndpoint 根据标识，获得对应endpoints
	SelectEndpoint(id string) (endpoint string, newSelect bool, err error)
	// 确认使用Endpoint，向Redis注册
	ConfirmEnpoint(endpoint, id string, checkIn bool) error
	ReleaseEndpoint(endpoint, id string)
	// 查询房间状态：endpoint: [id, ...]
	LoadStatus() map[string][]string
}

type LoadBalance int

const (
	// Random requests that connections are randomly distributed.
	Random LoadBalance = iota
	// RoundRobin requests that connections are distributed to a loop in a
	// round-robin fashion.
	RoundRobin
	// LeastConnections assigns the next accepted connection to the loop with
	// the least number of active connections.
	LeastConnections
)

type loop struct {
	endpoint string
	count    int // connection count
}

var (
	NoAvailableServiceEndPointError = errors.New("No available service endpoint")
	UpdateEndpointsInterval         = time.Second
)

type LBSelector struct {
	serviceName string

	serviceDiscovery Discovery
	loadbalance      LoadBalance

	loops     []*loop
	loopindex int
	loopslock sync.Mutex

	room *roomService
}

func NewSelector(serviceName string, config Config) Selector {
	mock := &LBSelector{
		serviceName:      serviceName,
		loadbalance:      LeastConnections,
		serviceDiscovery: &MockServiceDiscovery{},
		room:             NewRoomSerivce(config),
	}
	// 初次加载servicce的endpoint
	mock.refreshEndpoints()
	// todo: 服务越多，线程越多
	go mock.updateTrigger()
	return mock
}

func (m *LBSelector) ServiceName() string {
	return m.serviceName
}

func (m *LBSelector) SelectEndpoint(id string) (string, bool, error) {
	// 检查本地的endpoint缓存
	length := len(m.loops)
	if length == 0 {
		return "", false, NoAvailableServiceEndPointError
	}
	// 没有指定roomId，总是分配到第一个
	if id == DefaultRoomId {
		return m.loops[0].endpoint, false, nil
	}
	// 通过repository(redis)查询，如果id已和endpoint绑定
	if cachedEndpoint, ok := m.room.Where(id); ok {
		// 检查cache的endpoint和最新的endpoint是否一致
		for _, ep := range m.loops {
			if ep.endpoint == cachedEndpoint {
				return cachedEndpoint, false, nil
			}
		}
		// 清理错误的endpoint信息
		m.room.CheckOut(cachedEndpoint, id)
	}
	// todo: 如果updateTrigger将loops数量减少了，暂时不考虑
	var endpoint string
	if m.loadbalance == LeastConnections {
		minConnLoop := m.loops[0]
		for i := range m.loops {
			if m.loops[i].count < minConnLoop.count {
				minConnLoop = m.loops[i]
			}
		}
		minConnLoop.count++
		endpoint = minConnLoop.endpoint
	} else if m.loadbalance == RoundRobin {
		m.loopindex = (m.loopindex + 1) % length
		endpoint = m.loops[m.loopindex].endpoint
	}
	return endpoint, true, nil
}

// ConfirmEnpoint 确认
func (m *LBSelector) ConfirmEnpoint(endpoint, id string, checkIn bool) (err error) {
	if checkIn {
		// 写入room:endpoint
		// 覆盖写入endpoint:room = 1
		err = m.room.CheckIn(endpoint, id)
	} else {
		// 不写入room:endpoint
		// 写入endpoint:room += 1
		err = m.room.Enter(endpoint, id)
	}
	return
}

func (m *LBSelector) ReleaseEndpoint(endpoint, id string) {
	m.room.Leave(endpoint, id)
}

func (m *LBSelector) LoadStatus() map[string][]string {
	res := make(map[string][]string)
	for _, ep := range m.loops {
		// 重置本地loop数据，加载cache的
		roomInfo := m.room.Status(ep.endpoint)
		ep.count = len(roomInfo)
		res[ep.endpoint] = roomInfo
	}
	return res
}

// refreshEndpoints 覆盖更新endpoints信息
func (m *LBSelector) refreshEndpoints() {
	endpoints, err := m.serviceDiscovery.Endpoints(m.serviceName)
	if err != nil {
		log.Info("serviceDiscovery Endpoints error", err.Error())
		return
	}
	var newLoops []*loop
	for _, e := range endpoints {
		var exist bool
		for _, l := range m.loops {
			if l.endpoint == e {
				newLoops = append(newLoops, l)
				exist = true
				break
			}
		}
		if !exist {
			newLoops = append(newLoops, &loop{
				endpoint: e,
				count:    0,
			})
		}
	}
	m.loops = newLoops
}

// updateTrigger 开启定时器，定时查询最新endpoints
func (m *LBSelector) updateTrigger() {
	// 随机时间开启
	randStagger := time.Duration(uint64(rand.Int63()) % uint64(UpdateEndpointsInterval))
	<-time.After(randStagger)
	// 定时查询
	for {
		<-time.After(UpdateEndpointsInterval)
		m.refreshEndpoints()
	}
}
