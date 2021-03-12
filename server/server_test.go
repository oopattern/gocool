package server

import (
	"fmt"
	"log"
	"testing"
	"github.com/shima-park/agollo"
)

// 阿波罗配置中心: http://localhost:8070/
// 用户名: apollo
// 密码: admin

// consul UI页面
// http://localhost:8500/ui

var (
	route = RouteServer{}
	endpoint = ""
)

func init() {
	a, err := agollo.New("localhost:8080", "gongyi", agollo.AutoFetchOnCacheMiss())
	if err != nil {
		log.Fatalf("agollo config init failed: %+v", err)
	}
	port := a.Get("grpc_port", agollo.WithDefault("0"))
	endpoint = fmt.Sprintf("localhost:%s", port)
	ZapLogger.Info(fmt.Sprintf("grpc listen endpoint[%s]", endpoint))
}

func TestBuildGrpcServer(t *testing.T) {
	ZapLogger.Info(fmt.Sprintf("trace port[%d]", tracePort))
	s := NewServer(endpoint)
	s.RegisterService(route.RegisterServer)
	s.Run()
}