package main

import (
	"fmt"
	"github.com/shima-park/agollo"
	"github.com/oopattern/gocool/server"
	"github.com/oopattern/gocool/service"
	"github.com/oopattern/gocool/log"
	"github.com/oopattern/gocool/config"
)

var (
	route = service.RouteServer{}
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

func main() {
	s := server.NewServer(endpoint)
	s.RegisterService(route.RegisterServer)
	s.Run()
}
