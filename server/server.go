package main

// ref: https://github.com/apssouza22/grpc-production-go/blob/master/server/server.go
import (
	"net"
	"log"
	"google.golang.org/grpc"
)

func NewServer() (*grpc.Server) {
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	return grpcServer
}

func Run(listenAddr string, grpcServer *grpc.Server) {
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("failed to listen: +%v", err)
	}
	log.Fatal(grpcServer.Serve(lis))
}
