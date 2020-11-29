package server

// ref: https://github.com/apssouza22/grpc-production-go/blob/master/server/server.go
// ref: https://colobu.com/2017/04/17/dive-into-gRPC-interceptor/
// ref: https://www.liwenzhou.com/posts/Go/zap/
// ref: https://prometheus.io/
// ref: https://programmaticponderings.com/tag/prometheus/
// ref: https://medium.com/htc-research-engineering-blog/build-a-monitoring-dashboard-by-prometheus-grafana-741a7d949ec2
import (
	"fmt"
	"log"
	"net"
	"net/http"
	"google.golang.org/grpc"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	// grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
)

var (
	MetricsHttpPort = 9095
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
	// traceCfg := grpc_opentracing.UnaryServerInterceptor()
	logCfg := grpc_zap.UnaryServerInterceptor(ZapLogger)
	prometheusCfg := grpc_prometheus.UnaryServerInterceptor
	unaryOpt := grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(prometheusCfg, logCfg))

	var opts []grpc.ServerOption
	opts = append(opts, unaryOpt)
	s := grpc.NewServer(opts...)
	grpc_prometheus.Register(s)
	grpc_prometheus.EnableHandlingTimeHistogram()

	// Create a HTTP server for prometheus
	httpServer := &http.Server{Addr: fmt.Sprintf("0.0.0.0:%d", MetricsHttpPort)}
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal("Failed to start a http server")
		}
	}()

	// Create a TCP  server
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
