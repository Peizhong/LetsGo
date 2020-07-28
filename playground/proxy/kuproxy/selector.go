package main

import "errors"

// select upstream pod

type selector interface {
	// tellme 根据标识，获得对应endpoints
	// todo: 选取规则
	TellMe(id, serviceName string) (string, error)
}

type MockSelector struct {
	serviceDiscovery discovery
}

func NewSelector() selector {
	return &MockSelector{
		serviceDiscovery: MockServiceDiscovery{},
	}
}

var NoServiceEndPointError = errors.New("Cannot find service endPoint")

func (m *MockSelector) TellMe(id, serivceName string) (string, error) {
	endpoints, err := m.serviceDiscovery.Endpoints(serivceName)
	if err != nil {
		return "", err
	}
	if len(endpoints) == 0 {
		return "", NoServiceEndPointError
	}
	return endpoints[0], nil
}
