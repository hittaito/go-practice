package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

func myUnaryClientInterceptor(ctx context.Context, method string, req, res interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	fmt.Println("[unary] before interceptor")
	err := invoker(ctx, method, req, res, cc, opts...)
	fmt.Println("[unary] after interceptor")
	return err
}
