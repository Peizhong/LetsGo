package proxy

type Config interface {
	GetString(name string) string
}

const (
	redisAddrKey     = "RedisAddr"
	redisPasswordKey = "RedisPassword"
)

func NewConfig() Config {
	local := &localConfig{
		m: map[string]string{
			redisAddrKey:     "10.5.24.18:6379",
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
