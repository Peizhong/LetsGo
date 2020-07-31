package proxy

import "time"

type Config interface {
	GetString(name string) string
}

var (
	UpdateEndpointsInterval = 30 * time.Second

	// room:endpoint关联关系，一段时间后自动删除
	RoomExpireTime = time.Hour * 24

	redisAddrKey     = "RedisAddr"
	redisPasswordKey = "RedisPassword"

	// k8s中service的名字
	// todo：本地调试
	DefaultServiceName = "DefaultServiceName"
)

func NewConfig() Config {
	local := &localConfig{
		m: map[string]string{
			redisAddrKey:     "192.168.3.143:6379",
			redisPasswordKey: "",
		},
	}
	return local
}

type localConfig struct {
	m map[string]string
}

func (l *localConfig) GetString(name string) string {
	if v, ok := l.m[name]; ok {
		return v
	}
	return ""
}
