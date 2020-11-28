package server

// ref: https://github.com/apssouza22/grpc-production-go/blob/master/server/server.go
import (
	"net"
	"log"
	"google.golang.org/grpc"
)

type GrpcServer interface {
	Run()
	RegisterService(reg func(*grpc.Server))
}

type grpcServer struct {
	server *grpc.Server
	listener net.Listener
}

func (s *grpcServer) RegisterService(reg func(*grpc.Server)) {
	reg(s.server)
}

func (s *grpcServer) Run() {
	log.Fatal(s.server.Serve(s.listener))
}

func NewServer(endpoint string) GrpcServer {
	var opts []grpc.ServerOption
	s := grpc.NewServer(opts...)

	l, err := net.Listen("tcp", endpoint)
	if err != nil {
		log.Fatalf("failed to listen: +%v", err)
	}

	server := &grpcServer{
		server: s,
		listener: l,
	}
	return server
}
