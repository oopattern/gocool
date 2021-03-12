package server

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"github.com/shima-park/agollo"
	"github.com/oopattern/gocool/log"
	"github.com/oopattern/gocool/config"
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
	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGALRM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	s := NewServer(endpoint)
	s.RegisterService(route.RegisterServer)
	s.Run()
}