package helloworld

import (
	"flag"
	"net/http"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/peizhong/letsgo/framework/log"
	gw "github.com/peizhong/letsgo/services/grpc/pb/helloworld"
)

var (
	echoEndpoint = flag.String("echo_endpoint", "localhost:50051", "endpoint of YourService")
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := gw.RegisterGreeterHandlerFromEndpoint(ctx, mux, *echoEndpoint, opts)
	if err != nil {
		return err
	}

	return http.ListenAndServe(":8081", mux)
}

// https://github.com/grpc-ecosystem/grpc-gateway
// protoc -I C:/Users/wxyz/go/src  -I C:/Users/wxyz/go/src/github.com/googleapis/googleapis -I helloworld/ --go_out=plugins=grpc:helloworld helloworld/helloworld.proto
// generate a reverse proxy
// protoc -I C:/Users/wxyz/go/src  -I C:/Users/wxyz/go/src/github.com/googleapis/googleapis -I helloworld/ --grpc-gateway_out=logtostderr=true:helloworld helloworld/helloworld.proto
// protoc -I C:/Users/wxyz/go/src  -I C:/Users/wxyz/go/src/github.com/googleapis/googleapis -I helloworld/ --swagger_out=logtostderr=true:helloworld helloworld/helloworld.proto
// curl -X POST -k http://localhost:8080/v1/example/sayhello -d '{"name": "CoS is hname:"}'
func Run() {
	flag.Parse()
	defer glog.Flush()

	log.Info("api_helloworldproxy is on")
	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
