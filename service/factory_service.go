package service

import (
	"context"
	"google.golang.org/grpc"
	"github.com/oopattern/gocool/log"
	"github.com/oopattern/gocool/proto"
	"github.com/oopattern/gocool/server"
)

type FactoryServer struct {
	proto.UnimplementedFactoryServer
}

func (r *FactoryServer)RegisterServer(endpoint string, grpcServer *grpc.Server) {
	// register gRpc server
	proto.RegisterFactoryServer(grpcServer, r)
	// register gRpc gateway
	ctx := context.Background()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := proto.RegisterFactoryHandlerFromEndpoint(ctx, server.GatewayMux, endpoint, opts); err != nil {
		log.Fatal("register gateway err")
	}
}

func (r *RouteServer) CreateScheduler(ctx context.Context, req *proto.ConfigReq) (*proto.SchedulerResp, error) {
	log.Debug("rpc call ok")
	return &proto.SchedulerResp{
		Id: "default-scheduler",
	}, nil
}
