package main

import (
	"context"
	"io"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	tlog "github.com/opentracing/opentracing-go/log"
	"github.com/peizhong/letsgo/pkg/config"
	pb "github.com/peizhong/letsgo/playground/rpc/pb/twoway"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
)

// printFeature gets the feature for the given point.
func printFeature(client pb.TwoWayJobClient, point *pb.Point) {
	log.Printf("Getting feature for point (%d, %d)", point.Latitude, point.Longitude)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	feature, err := client.GetFeature(ctx, point)
	if err != nil {
		log.Fatalf("%v.GetFeatures(_) = _, %v: ", client, err)
	}
	log.Println(feature)
}

// printFeatures lists all the features within the given bounding Rectangle.
func printFeatures(client pb.TwoWayJobClient, rect *pb.Rectangle) {
	log.Printf("Looking for features within %v", rect)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.ListFeatures(ctx, rect)
	if err != nil {
		log.Fatalf("%v.ListFeatures(_) = _, %v", client, err)
	}
	for {
		feature, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.ListFeatures(_) = _, %v", client, err)
		}
		log.Println(feature)
	}
}

func randomPoint(r *rand.Rand) *pb.Point {
	lat := (r.Int31n(180) - 90) * 1e7
	long := (r.Int31n(360) - 180) * 1e7
	return &pb.Point{Latitude: lat, Longitude: long}
}

// runRecordRoute sends a sequence of points to server and expects to get a RouteSummary from server.
func runRecordRoute(client pb.TwoWayJobClient) {
	// Create a random number of random points
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	pointCount := int(r.Int31n(100)) + 2 // Traverse at least two points
	var points []*pb.Point
	for i := 0; i < pointCount; i++ {
		points = append(points, randomPoint(r))
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.RecordRoute(ctx)
	if err != nil {
		log.Fatalf("can not record: %v", err)
	}
	for _, point := range points {
		if err := stream.Send(point); err != nil {
			log.Fatalf("%v.Send(%v) = %v", stream, point, err)
		}
	}
	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
	}
	log.Printf("Route summary: %v", reply)
}

//OpenTracingClientInterceptor  rewrite client's interceptor with open tracing
func OpenTracingClientInterceptor(tracer opentracing.Tracer) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, resp interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		var parentCtx opentracing.SpanContext
		if parent := opentracing.SpanFromContext(ctx); parent != nil {
			parentCtx = parent.Context()
		}
		cliSpan := tracer.StartSpan(
			method,
			opentracing.ChildOf(parentCtx),
			ext.SpanKindRPCClient,
		)
		defer cliSpan.Finish()
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		} else {
			md = md.Copy()
		}
		err := tracer.Inject(cliSpan.Context(), opentracing.TextMap, md)
		if err != nil {
			grpclog.Errorf("inject to metadata err %v", err)
		}
		ctx = metadata.NewOutgoingContext(ctx, md)
		err = invoker(ctx, method, req, resp, cc, opts...)
		if err != nil {
			cliSpan.LogFields(tlog.String("err", err.Error()))
		}
		return err
	}
}

func main() {
	port := config.GrpcApp1Port
	addr := net.JoinHostPort("", strconv.Itoa(port))
	// 客户端只需要证书
	creds, _ := credentials.NewClientTLSFromFile(config.CertCrt, config.CertName)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds), grpc.WithBlock(),
		grpc.WithUnaryInterceptor(OpenTracingClientInterceptor(opentracing.GlobalTracer())))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	cli := pb.NewTwoWayJobClient(conn)
	printFeature(cli, &pb.Point{Latitude: 409146138, Longitude: -746188906})
	printFeature(cli, &pb.Point{Latitude: 0, Longitude: 0})
	printFeatures(cli, &pb.Rectangle{
		Lo: &pb.Point{Latitude: 400000000, Longitude: -750000000},
		Hi: &pb.Point{Latitude: 420000000, Longitude: -730000000},
	})
	runRecordRoute(cli)
}
