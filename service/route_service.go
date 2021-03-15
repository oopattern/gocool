package service

import (
	"context"
	"github.com/oopattern/gocool/log"
	"github.com/oopattern/gocool/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"strings"
)

type RouteServer struct {
	proto.UnimplementedObserveServer
}

func (r *RouteServer)RegisterServer(endpoint string, grpcServer *grpc.Server) {
	// register gRpc server
	proto.RegisterObserveServer(grpcServer, r)
}

func (r *RouteServer) SayRoute(ctx context.Context, req *proto.RouteReq) (*proto.RouteResp, error) {
	ip := "localhost"
	port := "0"
	if pr, ok := peer.FromContext(ctx); ok {
		addr := strings.Split(pr.Addr.String(), ":")
		if "[" != addr[0] {
			ip = addr[0]
		}
		port = addr[1]
	}
	log.Debug("rpc call ok")
	return &proto.RouteResp{
		Ip:   ip,
		Port: port,
	}, nil
}
