package server

// ref: https://github.com/apssouza22/grpc-production-go/blob/master/server/server.go
// ref: https://colobu.com/2017/04/17/dive-into-gRPC-interceptor/
// ref: https://www.liwenzhou.com/posts/Go/zap/
import (
	"net"
	"log"
	"google.golang.org/grpc"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
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
	unaryOpt := grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		grpc_zap.UnaryServerInterceptor(ZapLogger)))

	var opts []grpc.ServerOption
	opts = append(opts, unaryOpt)
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
