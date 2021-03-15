package server

import (
	"fmt"
	"github.com/oopattern/gocool/config"
	"github.com/oopattern/gocool/log"
	"github.com/oopattern/gocool/service"
	"github.com/shima-park/agollo"
	"testing"
)

var (
	factory = service.FactoryServer{}
	endpoint = ""
)

func init() {
	a, err := agollo.New(config.AgolloEndPoint, "gongyi", agollo.AutoFetchOnCacheMiss())
	if err != nil {
		log.Fatal("agollo config init failed: %+v", err)
	}
	port := a.Get("factory_port", agollo.WithDefault("0"))
	endpoint = fmt.Sprintf("localhost:%s", port)
	log.Info("grpc listen endpoint[%s]", endpoint)
}

func TestBuildGrpcServer(t *testing.T) {
	s := NewServer(endpoint)
	s.RegisterService(factory.RegisterServer)
	s.Run()
}