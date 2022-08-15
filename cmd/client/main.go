package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	mygrpc "github.com/hittaito/go-practice/pkg/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
func main() {
	scanner = bufio.NewScanner(os.Stdin)

	address := "localhost:8080"
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())

	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client = mygrpc.NewGreetingServiceClient(conn)

	for {
		fmt.Println("1: send requet")
		fmt.Println("2: exit")
		fmt.Print("please enter >")

		scanner.Scan()
		in := scanner.Text()

		switch in {
		case "1":
			Hello()
		case "2":
			goto M
		}
	}
M:
}
