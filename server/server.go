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
// ref: https://skyapm.github.io/document-cn-translation-of-skywalking/zh/8.0.0/
// ref: https://learn.hashicorp.com/tutorials/consul/get-started-service-discovery
// ref: http://blog.didispace.com/consul-service-discovery-exp/
// ref: https://juejin.cn/post/6844903794380111886
// ref: https://github.com/generals-space/gokit
// ref: https://stackoverflow.com/questions/30684262/different-ports-used-by-consul
// ref: https://www.consul.io/docs/install/ports
// ref: https://www.cnblogs.com/FireworksEasyCool/p/12782137.html
// ref: https://medium.com/swlh/rest-over-grpc-with-grpc-gateway-for-go-9584bfcbb835
// ref: https://jergoo.gitbooks.io/go-grpc-practice-guide/content/chapter3/gateway.html
// ref: https://github.com/jergoo/go-grpc-tutorial/tree/master/src/proto/google/api
import (
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
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
)

type GrpcServer interface {
	Run()
	RegisterService(reg func(*grpc.Server))
	GetListener() net.Listener
}

type grpcServer struct {
	server *grpc.Server
	http *http.Server
	listener net.Listener
}

func (s *grpcServer) GetListener() net.Listener {
	return s.listener
}

func (s *grpcServer) RegisterService(reg func(*grpc.Server)) {
	// register to gRpc
	reg(s.server)
	// register to consul
	endpoint := s.listener.Addr().String()
	for name, info := range s.server.GetServiceInfo() {
		if err := RegisterConsul(name, endpoint); err != nil {
			log.Fatalf("Failed to register service[%s]", name)
		}
		ZapLogger.Info(fmt.Sprintf("register service_name[%s], info[%+v]", name, info))
	}
	ZapLogger.Info(fmt.Sprintf("listen endpoint[%s]", endpoint))
}

func (s *grpcServer) Run() {
	log.Fatal(s.server.Serve(s.listener))
	// support REST gateway
	// log.Fatal(s.http.Serve(s.listener))
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

	// Create REST gateway
	h := RouteHttp(endpoint, s)

	// Create a TCP  server
	l, err := net.Listen("tcp", endpoint)
	if err != nil {
		log.Fatalf("failed to listen: +%v", err)
	}

	server := &grpcServer{
		server: s,
		http: h,
		listener: l,
	}
	return server
}
