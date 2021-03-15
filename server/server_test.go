package server

import (
	"fmt"
	"github.com/oopattern/gocool/config"
	"github.com/oopattern/gocool/log"
	"github.com/oopattern/gocool/service"
	"github.com/shima-park/agollo"
	"testing"
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
	a, err := agollo.New(config.AgolloEndPoint, "gongyi", agollo.AutoFetchOnCacheMiss())
	if err != nil {
		log.Fatal("agollo config init failed: %+v", err)
	}
	port := a.Get("grpc_port", agollo.WithDefault("0"))
	endpoint = fmt.Sprintf("localhost:%s", port)
	log.Info("grpc listen endpoint[%s]", endpoint)
}

func TestBuildGrpcServer(t *testing.T) {
	s := NewServer(endpoint)
	s.RegisterService(route.RegisterServer)
	s.Run()
}