package server

import (
	"fmt"
	"context"
	"google.golang.org/grpc"
	"github.com/oopattern/gocool/proto"
)

type routeServer struct {
	proto.UnimplementedObserveServer
}

func (r *routeServer)RegisterServer(grpcServer *grpc.Server) {
	proto.RegisterObserveServer(grpcServer, r)
}

func (r *routeServer) SayRoute(ctx context.Context, req *proto.RouteReq) (*proto.RouteResp, error) {
	fmt.Println(req.GetName())
	return &proto.RouteResp{
		Ip:   "localhost",
		Port: fmt.Sprintf("%d", 7777),
	}, nil
}
