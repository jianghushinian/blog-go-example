// Package main implements a server for Greeter service.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	genericapiserver "k8s.io/apiserver/pkg/server"
	// genericapiserver "github.com/jianghushinian/blog-go-example/gracefulstop/pkg/server"

	"github.com/jianghushinian/blog-go-example/gracefulstop/grpc/pb"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())

	duration, _ := time.ParseDuration(in.GetDuration())
	time.Sleep(duration)

	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	<-genericapiserver.SetupSignalHandler()
	log.Printf("Shutdown Server...")
	s.GracefulStop()
	log.Println("gRPC server graceful shutdown completed")
}
