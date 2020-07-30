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
	// SelectEndpoint 根据标识，获得对应endpoints，并分配房间。
	// 如果后续操作失败，通过ReleaseEndpoint释放房间
	SelectEndpoint(id string) (string, error)
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
	SelectServiceEndPointFailed     = errors.New("Select service endPoint failed")
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
	mock.refreshEndpoints()
	// stopCh := make(chan struct{})
	go mock.updateTrigger()
	return mock
}

func (m *LBSelector) SelectEndpoint(id string) (string, error) {
	// 通过repository(redis)查询，如果id已和endpoint绑定
	if cachedEndpoint, ok := m.room.Where(id); ok {
		// 检查cache的endpoint和最新的endpoint是否一致
		for _, ep := range m.loops {
			if ep.endpoint == cachedEndpoint {
				// 房间再增加一个连接
				m.room.Enter(cachedEndpoint, id)
				return cachedEndpoint, nil
			}
		}
		// 清理错误的endpoint信息
		m.room.CheckOut(cachedEndpoint, id)
	}
	length := len(m.loops)
	if length == 0 {
		return "", NoAvailableServiceEndPointError
	}
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
	if endpoint != "" {
		m.room.CheckIn(endpoint, id)
		return endpoint, nil
	}
	return "", SelectServiceEndPointFailed
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
