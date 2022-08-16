package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

func myUnaryServerInterceptor1(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, hander grpc.UnaryHandler) (interface{}, error) {
	log.Println("[unary] before interceptor")
	res, err := hander(ctx, req)
	log.Println("[unary] after request")
	return res, err
}
