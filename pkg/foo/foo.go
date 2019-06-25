package foo

import (
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/transport"
)

// mockgen -source=foo.go --destination=../mock_foo/foo.go

type Foo interface {
	Bar(x int) int
	Get()
	Set(x int)
	DoSomething(int, string) error
}

type Registry interface {
	Init(...registry.Option) error
	Options() registry.Options
	Register(*registry.Service, ...registry.RegisterOption) error
	Deregister(*registry.Service) error
	GetService(string) ([]*registry.Service, error)
	ListServices() ([]*registry.Service, error)
	Watch(...registry.WatchOption) (registry.Watcher, error)
	String() string
}

type Transport interface {
	Init(...transport.Option) error
	Options() transport.Options
	Dial(addr string, opts ...transport.DialOption) (transport.Client, error)
	Listen(addr string, opts ...transport.ListenOption) (transport.Listener, error)
	String() string
}

type Broker interface {
	Options() broker.Options
	Address() string
	Connect() error
	Disconnect() error
	Init(...broker.Option) error
	Publish(string, *broker.Message, ...broker.PublishOption) error
	Subscribe(string, broker.Handler, ...broker.SubscribeOption) (broker.Subscriber, error)
	String() string
}

func f() {
	var x = [5]int{1, 2, 3, 4, 5}
	_ = x
}
