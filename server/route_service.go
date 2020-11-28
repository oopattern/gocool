package server

import (
	"fmt"
	"strings"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"github.com/oopattern/gocool/proto"
)

type RouteServer struct {
	proto.UnimplementedObserveServer
}

func (r *RouteServer)RegisterServer(grpcServer *grpc.Server) {
	proto.RegisterObserveServer(grpcServer, r)
}

func (r *RouteServer) SayRoute(ctx context.Context, req *proto.RouteReq) (*proto.RouteResp, error) {
	fmt.Println(req.GetName())
	ip := "localhost"
	port := "0"
	if pr, ok := peer.FromContext(ctx); ok {
		 addr := strings.Split(pr.Addr.String(), ":")
		 if "[" != addr[0] {
			ip = addr[0]
		 }
		 port = addr[1]
	}
	return &proto.RouteResp{
		Ip:   ip,
		Port: port,
	}, nil
}
