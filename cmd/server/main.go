package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	mygrpc "github.com/hittaito/go-practice/pkg/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type myServer struct {
	mygrpc.UnimplementedGreetingServiceServer
}

func (s *myServer) Hello(ctx context.Context, req *mygrpc.HelloRequest) (*mygrpc.HelloResponse, error) {
	return &mygrpc.HelloResponse{
		Message: fmt.Sprintf("Hello, %s!", req.GetName()),
	}, nil
}
func (s *myServer) HelloServerStream(req *mygrpc.HelloRequest, stream mygrpc.GreetingService_HelloServerStreamServer) error {
	resCount := 5
	for i := 0; i < resCount; i++ {
		err := stream.Send(&mygrpc.HelloResponse{
			Message: fmt.Sprintf("[%d] Hello %s", i, req.GetName()),
		})
		if err != nil {
			return err
		}
		time.Sleep(time.Second)
	}
	return nil
}

func NewMyServer() *myServer {
	return &myServer{}
}

func main() {
	port := 8080
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()

	mygrpc.RegisterGreetingServiceServer(server, NewMyServer())

	reflection.Register(server)

	go func() {
		log.Printf("start grpc port %v", port)
		server.Serve(listener)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	log.Println("stop server")
	server.GracefulStop()
}
