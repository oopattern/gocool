package service

import (
	"context"
	"github.com/oopattern/gocool/log"
	"github.com/oopattern/gocool/proto"
	"google.golang.org/grpc"
)

type FactoryServer struct {
	proto.UnimplementedFactoryServer
}

func (r *FactoryServer)RegisterServer(endpoint string, grpcServer *grpc.Server) {
	// register gRpc server
	proto.RegisterFactoryServer(grpcServer, r)
}

func (r *RouteServer) CreateScheduler(ctx context.Context, req *proto.ConfigReq) (*proto.SchedulerResp, error) {
	log.Debug("rpc call ok")
	return &proto.SchedulerResp{
		Id: "default-scheduler",
	}, nil
}