package middlewares

import (
	"time"

	"github.com/peizhong/letsgo/framework"
	"github.com/peizhong/letsgo/gonet"
)

type CacheMiddleware struct {
}

func (m *CacheMiddleware) Invoke(c *gonet.Context, ch chan<- struct{}, next func(*gonet.Context, chan<- struct{}) error) (err error) {
	c.SayHi("最后一个中间件，没有next")
	if value, exist := cacheRepo.TryGet(c.SrcPath); exist {
		c.Response = value.Response
	}
	framework.SayHiOrm()
	if ch != nil {
		close(ch)
	}
	return nil
}

type CachedGatewayResponse struct {
	Response    *gonet.GatewayResponse
	Exipre      time.Time `remark:"过期时间“`
	RequestTime int64     `remark:"请求次数"`
}

type CacheRepository struct {
	cache map[string]*CachedGatewayResponse
}

func (repo CacheRepository) TryGet(key string) (*CachedGatewayResponse, bool) {
	if value, exist := repo.cache[key]; exist {
		if time.Since(value.Exipre) < 0 {
			// todo: renew expire
			return value, true
		}
		delete(repo.cache, key)
	}
	return nil, false
}

func (repo CacheRepository) TryAdd(key string, response *CachedGatewayResponse) bool {
	if _, exist := repo.cache[key]; !exist {
		repo.cache[key] = response
	}
	return true
}

var cacheRepo = CacheRepository{
	cache: make(map[string]*CachedGatewayResponse),
}
