package server

// ref: https://github.com/apssouza22/grpc-production-go/blob/master/server/server.go
// ref: https://colobu.com/2017/04/17/dive-into-gRPC-interceptor/
// ref: https://www.liwenzhou.com/posts/Go/zap/
// ref: https://prometheus.io/
// ref: https://programmaticponderings.com/tag/prometheus/
// ref: https://programmaticponderings.com/tag/jaeger/
// ref: https://medium.com/htc-research-engineering-blog/build-a-monitoring-dashboard-by-prometheus-grafana-741a7d949ec2
// ref: https://jishuin.proginn.com/p/763bfbd310b3
// ref: https://medium.com/opentracing/tracing-http-request-latency-in-go-with-opentracing-7cc1282a100a
// ref: https://www.selinux.tech/golang/grpc/grpc-tracing
// ref: https://github.com/bigbully/Dapper-translation
// ref: https://juejin.cn/post/6871928187123826702
// ref: https://blog.csdn.net/zhounixing/article/details/105815910
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
	// jaeger "github.com/uber/jaeger-client-go"
	// opentracing "github.com/opentracing/opentracing-go"
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
