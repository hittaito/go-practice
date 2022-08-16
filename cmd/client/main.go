package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	mygrpc "github.com/hittaito/go-practice/pkg/grpc"
	_ "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var (
	scanner *bufio.Scanner
	client  mygrpc.GreetingServiceClient
)

func Hello() {
	fmt.Print("enter your name >")
	scanner.Scan()

	name := scanner.Text()

	req := &mygrpc.HelloRequest{
		Name: name,
	}
	res, err := client.Hello(context.Background(), req)
	if err != nil {
		panic(err)
	}
	fmt.Println(res.GetMessage())
}
func HelloStream() {
	fmt.Print("enter your name >")
	scanner.Scan()

	name := scanner.Text()

	req := &mygrpc.HelloRequest{
		Name: name,
	}
	stream, err := client.HelloServerStream(context.Background(), req)
	if err != nil {
		panic(err)
	}
	for {
		res, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("all response have received")
			break
		}
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(res)
	}
}
func HelloClientStream() {
	stream, err := client.HelloClientStream(context.Background())
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("enter your name")
	for {
		scanner.Scan()
		name := scanner.Text()

		if name == "" {
			break
		}
		err = stream.Send(&mygrpc.HelloRequest{
			Name: name,
		})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(res.GetMessage())
}
func HelloBiStream() {
	stream, err := client.HelloBiStream(context.Background())
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("enter your name")

	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("finish stream")
				close(waitc)
				return
			}
			if err != nil {
				return
			}
			log.Printf("got message: %s", in.GetMessage())
		}
	}()

	for {
		scanner.Scan()
		name := scanner.Text()

		if name == "" {
			break
		}
		err = stream.Send(&mygrpc.HelloRequest{
			Name: name,
		})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
	err = stream.CloseSend()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	<-waitc
}

func FailHello() {

	req := &mygrpc.HelloRequest{
		Name: "hoge",
	}
	res, err := client.FailHello(context.Background(), req)
	if err != nil {
		fmt.Println(err)
		if stat, ok := status.FromError(err); ok {
			fmt.Printf("code: %s\n", stat.Code())
			fmt.Printf("message: %s\n", stat.Message())
			fmt.Printf("detail: %s\n", stat.Details()...)
		}
	} else {
		fmt.Println(res.GetMessage())
	}
}
func main() {
	scanner = bufio.NewScanner(os.Stdin)

	address := "localhost:8080"
	conn, err := grpc.Dial(
		address,
		grpc.WithStreamInterceptor(myStreamClientInterceptor),
		grpc.WithUnaryInterceptor(myUnaryClientInterceptor),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)

	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client = mygrpc.NewGreetingServiceClient(conn)

	for {
		fmt.Println("")
		fmt.Println("1: send requet")
		fmt.Println("2: send server stream request")
		fmt.Println("3: send client stream request")
		fmt.Println("4: send bi stream request")
		fmt.Println("5: server side error request")
		fmt.Println("6: exit")
		fmt.Print("please enter >")

		scanner.Scan()
		in := scanner.Text()

		switch in {
		case "1":
			Hello()
		case "2":
			HelloStream()
		case "3":
			HelloClientStream()
		case "4":
			HelloBiStream()
		case "5":
			FailHello()
		case "6":
			goto M
		}
	}
M:
}
