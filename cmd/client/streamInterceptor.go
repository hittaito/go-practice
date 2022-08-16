package main

import (
	"context"
	"errors"
	"io"
	"log"

	"google.golang.org/grpc"
)

type myClientStreamWrapper struct {
	grpc.ClientStream
}

func (s *myClientStreamWrapper) SendMsg(m interface{}) error {
	log.Println("[stream] send message")
	return s.ClientStream.SendMsg(m)
}
func (s *myClientStreamWrapper) RecvMsg(m interface{}) error {
	err := s.ClientStream.RecvMsg(m)

	if !errors.Is(err, io.EOF) {
		log.Println("[stream] recv message")
	}
	return err
}
func (s *myClientStreamWrapper) CloseSend() error {
	err := s.ClientStream.CloseSend()
	log.Println("[stream] close send")
	return err
}

func myStreamClientInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	log.Printf("[stream] before interceptor (%s)\n", method)
	stream, err := streamer(ctx, desc, cc, method, opts...)

	return &myClientStreamWrapper{stream}, err
}
