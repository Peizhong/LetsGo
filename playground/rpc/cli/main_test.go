package main

import (
	"github.com/gogo/protobuf/proto"
	"github.com/peizhong/letsgo/playground/rpc/pb/helloworld"
	"testing"
)

func TestMarshall(t *testing.T) {
	request := &helloworld.HelloRequest{Name: "www"}
	data, err := proto.Marshal(request)
	if err != nil {
		return
	}
	newRequest := &helloworld.HelloRequest{}
	err = proto.Unmarshal(data, newRequest)
	if err != nil {
		return
	}
}
