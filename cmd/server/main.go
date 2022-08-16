package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	mygrpc "github.com/hittaito/go-practice/pkg/grpc"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type myServer struct {
	mygrpc.UnimplementedGreetingServiceServer
}

func (s *myServer) Hello(ctx context.Context, req *mygrpc.HelloRequest) (*mygrpc.HelloResponse, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		log.Println(md)
	}

	headerMD := metadata.New(map[string]string{"type": "unary", "from": "server", "in": "header"})
	if err := grpc.SetHeader(ctx, headerMD); err != nil {
		return nil, err
	}

	trailerMD := metadata.New(map[string]string{"type": "unary", "from": "server", "in": "trailer"})
	if err := grpc.SetTrailer(ctx, trailerMD); err != nil {
		return nil, err
	}
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
func (s *myServer) HelloClientStream(stream mygrpc.GreetingService_HelloClientStreamServer) error {
	nameList := make([]string, 0)

	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			message := fmt.Sprintf("hello, %v", nameList)
			err = stream.SendAndClose(&mygrpc.HelloResponse{
				Message: message,
			})
			return err
		}
		if err != nil {
			return err
		}
		nameList = append(nameList, req.GetName())
	}
}
func (s *myServer) HelloBiStream(stream mygrpc.GreetingService_HelloBiStreamServer) error {
	if md, ok := metadata.FromIncomingContext(stream.Context()); ok {
		log.Println(md)
	}

	headerMD := metadata.New(map[string]string{"type": "stream", "from": "server", "in": "header"})
	if err := stream.SetHeader(headerMD); err != nil {
		return err
	}
	trailerMD := metadata.New(map[string]string{"type": "stream", "from": "server", "in": "trailer"})
	stream.SetTrailer(trailerMD)
	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}

		message := fmt.Sprintf("Hello, %v", req.GetName())

		if err := stream.Send(&mygrpc.HelloResponse{
			Message: message,
		}); err != nil {
			return err
		}
	}
}
func (s *myServer) FailHello(ctx context.Context, req *mygrpc.HelloRequest) (*mygrpc.HelloResponse, error) {
	stat := status.New(codes.Unknown, "unknown error occurred")
	stat, _ = stat.WithDetails(&errdetails.DebugInfo{
		Detail: "detail reason here",
	})
	err := stat.Err()
	return nil, err
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

	server := grpc.NewServer(
		grpc.UnaryInterceptor(myUnaryServerInterceptor1),
		grpc.StreamInterceptor(myStreamServerInterceptor),
	)

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
