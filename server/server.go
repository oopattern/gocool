package server

import (
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	// jaeger "github.com/uber/jaeger-client-go"
	// opentracing "github.com/opentracing/opentracing-go"
)

var (
	MetricsHttpPort = 9095
	GatewayPort = 8006
	GatewayMux = runtime.NewServeMux()
)

type GrpcServer interface {
	Run()
	RegisterService(reg func(endpoint string, server *grpc.Server))
	GetListener() net.Listener
}

type grpcServer struct {
	server *grpc.Server
	listener net.Listener
}

func (s *grpcServer) GetListener() net.Listener {
	return s.listener
}

func (s *grpcServer) RegisterService(reg func(endpoint string, server *grpc.Server)) {
	endpoint := s.listener.Addr().String()
	// register to gRpc
	reg(endpoint, s.server)
	// register to consul
	for name, info := range s.server.GetServiceInfo() {
		if err := RegisterConsul(name, endpoint); err != nil {
			log.Fatalf("Failed to register service[%s]", name)
		}
		ZapLogger.Info(fmt.Sprintf("register service_name[%s], info[%+v]", name, info))
	}
	ZapLogger.Info(fmt.Sprintf("listen endpoint[%s]", endpoint))
}

func (s *grpcServer) Run() {
	// run gRpc gateway
	StartGateway(fmt.Sprintf(":%d", GatewayPort), s.server)
	// run gRpc server
	log.Fatal(s.server.Serve(s.listener))
}

func NewServer(endpoint string) GrpcServer {
	logCfg := grpc_zap.UnaryServerInterceptor(ZapLogger)
	prometheusCfg := grpc_prometheus.UnaryServerInterceptor
	unaryOpt := grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(prometheusCfg, logCfg))

	var opts []grpc.ServerOption
	opts = append(opts, unaryOpt)
	s := grpc.NewServer(opts...)

	// Create a HTTP server for prometheus
	grpc_prometheus.Register(s)
	grpc_prometheus.EnableHandlingTimeHistogram()
	prometheusServer := &http.Server{Addr: fmt.Sprintf("0.0.0.0:%d", MetricsHttpPort)}
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		if err := prometheusServer.ListenAndServe(); err != nil {
			log.Fatal("Failed to start a http server")
		}
	}()

	// Create a TCP  server
	l, err := net.Listen("tcp", endpoint)
	if err != nil {
		log.Fatalf("failed to listen: %+v", err)
	}

	server := &grpcServer{
		server: s,
		listener: l,
	}
	return server
}
