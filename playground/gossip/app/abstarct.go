package app

const (
	ActionNone = 1 << iota
	ActionAdd
	ActionDel
)

type MemberInfo struct {
	Name   string
	Addr   string
	Status int
}

type ICloud interface {
	Save(path string) error
	Load(path string) error

	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
	Delete(key string) error

	Info() ([]*MemberInfo, error)

	Benchmark(n int) error
}

type dataType struct {
	Id    []byte
	Value interface{}
	Time  int64
}

type dataStore map[string]dataType

// kv变更的记录，附载在broadcast中传输
type update struct {
	Action int // add, del
	Data   dataStore
}
