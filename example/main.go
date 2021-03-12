package main

import (
	/*
	"os"
	"os/signal"
	"syscall"
	*/
	"fmt"
	"github.com/shima-park/agollo"
	"github.com/oopattern/gocool/server"
	"github.com/oopattern/gocool/log"
)

// 阿波罗配置中心: http://localhost:8070/
// 用户名: apollo
// 密码: admin

// consul UI页面
// http://localhost:8500/ui

var (
	route = server.RouteServer{}
	endpoint = ""
)

func init() {
	a, err := agollo.New("localhost:8080", "gongyi", agollo.AutoFetchOnCacheMiss())
	if err != nil {
		log.Fatal("agollo config init failed: %+v", err)
	}
	port := a.Get("grpc_port", agollo.WithDefault("0"))
	endpoint = fmt.Sprintf("localhost:%s", port)
	log.Info("grpc listen endpoint[%s]", endpoint)
}

func main() {
	/*
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGALRM)
		log.Error("catch signal[%s], process is ready to quit", <-c)
		os.Exit(0)
	}()
	*/
	s := server.NewServer(endpoint)
	s.RegisterService(route.RegisterServer)
	s.Run()
}
