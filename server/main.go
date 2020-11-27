package main

import (
	"fmt"
	"log"
	"context"
	"net"
	"google.golang.org/grpc"
	"gocool/proto"
)

type routeServer struct {
	proto.UnimplementedObserveServer
}

var (
	RouteHost = "localhost"
	RoutePort = 7777
)

func (r *routeServer) SayRoute(ctx context.Context, req *proto.RouteReq) (*proto.RouteResp, error) {
	fmt.Println(req.GetName())
	return &proto.RouteResp{
		Ip:   RouteHost,
		Port: fmt.Sprintf("%d", RoutePort),
	}, nil
}

func main() {
	fmt.Println("hello world")
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", RouteHost, RoutePort))
	if err != nil {
		log.Fatalf("failed to listen: +%v", err)
	}
	r := routeServer{}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	proto.RegisterObserveServer(grpcServer, &r)
	log.Fatal(grpcServer.Serve(lis))
}
