package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	consul_api "github.com/hashicorp/consul/api"
	"github.com/peizhong/letsgo/pkg/config"
	"github.com/peizhong/letsgo/playground/rpc/pb/helloworld"
	pb "github.com/peizhong/letsgo/playground/rpc/pb/twoway"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type server struct {
	savedFeatures []*pb.Feature // read-only after initialized
	mu            sync.Mutex    // protects routeNotes
	routeNotes    map[string][]*pb.RouteNote
}

func (s *server) Simple(context.Context, *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{Message: "hello"}, nil
}

func (s *server) GetFeature(ctx context.Context, point *pb.Point) (*pb.Feature, error) {
	for _, feature := range s.savedFeatures {
		if proto.Equal(feature.Location, point) {
			return feature, nil
		}
	}
	// No feature was found, return an unnamed feature
	return &pb.Feature{Location: point}, nil
}

func (s *server) ListFeatures(rect *pb.Rectangle, stream pb.TwoWayJob_ListFeaturesServer) error {
	inRange := func(point *pb.Point, rect *pb.Rectangle) bool {
		left := math.Min(float64(rect.Lo.Longitude), float64(rect.Hi.Longitude))
		right := math.Max(float64(rect.Lo.Longitude), float64(rect.Hi.Longitude))
		top := math.Max(float64(rect.Lo.Latitude), float64(rect.Hi.Latitude))
		bottom := math.Min(float64(rect.Lo.Latitude), float64(rect.Hi.Latitude))

		if float64(point.Longitude) >= left &&
			float64(point.Longitude) <= right &&
			float64(point.Latitude) >= bottom &&
			float64(point.Latitude) <= top {
			return true
		}
		return false
	}
	for _, feature := range s.savedFeatures {
		if inRange(feature.Location, rect) {
			if err := stream.Send(feature); err != nil {
				return err
			}
		}
	}
	return nil
}

func calcDistance(p1 *pb.Point, p2 *pb.Point) int32 {
	toRadians := func(num float64) float64 {
		return num * math.Pi / float64(180)
	}
	const CordFactor float64 = 1e7
	const R = float64(6371000) // earth radius in metres
	lat1 := toRadians(float64(p1.Latitude) / CordFactor)
	lat2 := toRadians(float64(p2.Latitude) / CordFactor)
	lng1 := toRadians(float64(p1.Longitude) / CordFactor)
	lng2 := toRadians(float64(p2.Longitude) / CordFactor)
	dlat := lat2 - lat1
	dlng := lng2 - lng1

	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(dlng/2)*math.Sin(dlng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := R * c
	return int32(distance)
}

func (s *server) RecordRoute(stream pb.TwoWayJob_RecordRouteServer) error {
	var pointCount, featureCount, distance int32
	var lastPoint *pb.Point
	startTime := time.Now()
	for {
		point, err := stream.Recv()
		if err == io.EOF {
			endTime := time.Now()
			return stream.SendAndClose(&pb.RouteSummary{
				PointCount:   pointCount,
				FeatureCount: featureCount,
				Distance:     distance,
				ElapsedTime:  int32(endTime.Sub(startTime).Seconds()),
			})
		}
		if err != nil {
			return err
		}
		pointCount++
		for _, feature := range s.savedFeatures {
			if proto.Equal(feature.Location, point) {
				featureCount++
			}
		}
		if lastPoint != nil {
			distance += calcDistance(lastPoint, point)
		}
		lastPoint = point
	}
}

func serialize(point *pb.Point) string {
	return fmt.Sprintf("%d %d", point.Latitude, point.Longitude)
}

func (s *server) RouteChat(stream pb.TwoWayJob_RouteChatServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		key := serialize(in.Location)
		s.mu.Lock()
		s.routeNotes[key] = append(s.routeNotes[key], in)
		rn := make([]*pb.RouteNote, len(s.routeNotes[key]))
		copy(rn, s.routeNotes[key])
		s.mu.Unlock()

		for _, note := range rn {
			if err := stream.Send(note); err != nil {
				return err
			}
		}
	}
}

func (s *server) loadFeatures(filePath string) {
	var data []byte
	var err error
	data, err = ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to load default features: %v", err)
	}

	if err := json.Unmarshal(data, &s.savedFeatures); err != nil {
		log.Fatalf("Failed to load default features: %v", err)
	}
}

func (s *server) Check(context.Context, *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	resp := &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}
	return resp, nil
}

func (s *server) Watch(req *grpc_health_v1.HealthCheckRequest, srv grpc_health_v1.Health_WatchServer) error {
	if req.Service == "" {
		return fmt.Errorf("unknow service")
	}
	srv.Send(&grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	})
	mock := make(chan int32)
	for {
		s := <-mock
		srv.Send(&grpc_health_v1.HealthCheckResponse{
			Status: grpc_health_v1.HealthCheckResponse_ServingStatus(s),
		})
	}
}

func main() {
	port := config.GrpcApp1Port
	addr := net.JoinHostPort("", strconv.Itoa(port))
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("listening at", lis.Addr().String())
	// tls 服务器需要密钥和证书
	creds, err := credentials.NewServerTLSFromFile(config.CertCrt, config.CertKey)
	if err != nil {
		log.Fatalf("failed to set creds: %v", err)
	}
	// tls
	s := grpc.NewServer(grpc.Creds(creds))
	// 暂时不用tls，给consul用
	s = grpc.NewServer(
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
		grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor))
	gServer := &server{}
	gServer.loadFeatures(filepath.Join(config.WorkspaceDir, "playground/conf/example.json"))
	pb.RegisterTwoWayJobServer(s, gServer)
	// 注册 健康检查
	grpc_health_v1.RegisterHealthServer(s, gServer)
	// 注册 反射
	// ./grpcurl -insecure localhost:8091 list
	// ./grpcurl -insecure localhost:8091 describe grpc.health.v1.Health.Watch
	reflection.Register(s)
	// consul 注册
	// 如果grpc启用了tls，consul也要配
	// ./consul agent -dev -ui
	if true {
		consul, err := consul_api.NewClient(consul_api.DefaultConfig())
		if err != nil {
			panic(err)
		}
		members, err := consul.Agent().Members(true)
		if err != nil {
			panic(err)
		}
		for _, m := range members {
			log.Println(m.Name, m.Addr)
		}
		err = consul.Agent().ServiceRegister(&consul_api.AgentServiceRegistration{
			ID:      lis.Addr().String(),
			Name:    "twoway.TwoWayJob",
			Address: lis.Addr().String(),
			Port:    port,
			Check: &consul_api.AgentServiceCheck{
				Interval:   "1000ms",
				GRPC:       fmt.Sprintf("%v/grpc.reflection.v1alpha.ServerReflection/Check", lis.Addr().String()),
				GRPCUseTLS: false,
			},
		})
		if err != nil {
			log.Println(err.Error())
		}
	}
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
