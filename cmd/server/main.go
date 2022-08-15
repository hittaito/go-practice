package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

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
