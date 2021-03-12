package server

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/oopattern/gocool/config"
	"github.com/oopattern/gocool/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
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
			log.Fatal("Failed to register service[%s]", name)
		}
		log.Info("register service_name[%s], info[%+v]", name, info)
	}
}

func (s *grpcServer) Run() {
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGALRM)
		log.Error("catch signal[%s], process is ready to quit", <-c)
		os.Exit(0)
	}()

	// run gRpc gateway
	StartGateway(config.GatewayEndPoint, s.server)
	// run gRpc server
	if err := s.server.Serve(s.listener); err != nil {
		log.Error("server catch signal to quit")
	}
}

func NewServer(endpoint string) GrpcServer {
	logCfg := grpc_zap.UnaryServerInterceptor(log.ZapLogger)
	prometheusCfg := grpc_prometheus.UnaryServerInterceptor
	unaryOpt := grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(prometheusCfg, logCfg))

	var opts []grpc.ServerOption
	opts = append(opts, unaryOpt)
	s := grpc.NewServer(opts...)

	// Create a HTTP server for prometheus
	grpc_prometheus.Register(s)
	grpc_prometheus.EnableHandlingTimeHistogram()
	prometheusServer := &http.Server{Addr: config.MetricsEndPoint}
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		if err := prometheusServer.ListenAndServe(); err != nil {
			log.Fatal("Failed to start a http server")
		}
	}()

	// Create a TCP  server
	l, err := net.Listen("tcp", endpoint)
	if err != nil {
		log.Fatal("failed to listen: %+v", err)
	}

	server := &grpcServer{
		server: s,
		listener: l,
	}
	return server
}
