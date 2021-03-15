package server

import (
	"strings"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"github.com/oopattern/gocool/log"
	"github.com/oopattern/gocool/proto"
)

type RouteServer struct {
	proto.UnimplementedObserveServer
}

func (r *RouteServer)RegisterServer(endpoint string, grpcServer *grpc.Server) {
	// register gRpc server
	proto.RegisterObserveServer(grpcServer, r)
	// register gRpc gateway
	ctx := context.Background()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := proto.RegisterObserveHandlerFromEndpoint(ctx, GatewayMux, endpoint, opts); err != nil {
		log.Fatal("register gateway err")
	}
}

func (r *RouteServer) Looks() error {
	return nil
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
