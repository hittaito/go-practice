package main

import (
	"errors"
	"io"
	"log"

	"google.golang.org/grpc"
)

type myServerStreamWrapper1 struct {
	grpc.ServerStream
}

func (s *myServerStreamWrapper1) RecvMsg(m interface{}) error {
	err := s.ServerStream.RecvMsg(m)

	if !errors.Is(err, io.EOF) {
		log.Println("[server stream] receive message")
	}
	return err
}
func (s *myServerStreamWrapper1) SendMsg(m interface{}) error {
	log.Println("[server stream] send message")
	return s.ServerStream.SendMsg(m)
}

func myStreamServerInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	log.Println("[server stream] before interceptor")
	err := handler(srv, &myServerStreamWrapper1{ss})
	log.Println("[server stream] after interceptor")
	return err
}
