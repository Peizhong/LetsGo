package proxy

import (
	"context"
	"fmt"
	"log"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

//go:generate mockgen -destination ./mock/room.go -source discovery.go
type Discovery interface {
	// endpoints 返回endpoints信息
	Endpoints(serviceName string) ([]string, error)
}

// https://github.com/kubernetes/client-go/tree/master/examples

type K8sServiceDiscovery struct {
}

func (*K8sServiceDiscovery) Endpoints(serviceName string) ([]string, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		if err == rest.ErrNotInCluster {
			log.Println("not in cluster")
			// microk8s kubectl proxy --accept-hosts=.* --address=0.0.0.0
			config = &rest.Config{
				Host: "192.168.3.143:8001",
			}
		} else {
			return nil, err
		}
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	// microk8s.kubectl proxy --accept-hosts=.* --address=0.0.0.0
	// http://192.168.3.143:8001/api/v1/namespaces/default/endpoints/nginx-service
	endpoint, err := client.CoreV1().Endpoints("default").Get(context.Background(), serviceName, v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	res := []string{}
	for _, ep := range endpoint.Subsets {
		for _, addr := range ep.Addresses {
			// todo: 如果pod暴露了多个端口，subsets.ports会有多个，name区分
			res = append(res, fmt.Sprintf("%s:%d", addr.IP, ep.Ports[0].Port))
		}
	}
	return res, nil
}

type MockServiceDiscovery struct {
}

func (*MockServiceDiscovery) Endpoints(serviceName string) ([]string, error) {
	return []string{"localhost:3000", "localhost:3001"}, nil
}
