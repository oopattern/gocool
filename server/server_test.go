package server

import (
	"fmt"
	"testing"
)

var (
	route = RouteServer{}
	endpoint = fmt.Sprintf("%s:%d", "localhost", 7777)
)

func TestBuildGrpcServer(t *testing.T) {
	ZapLogger.Debug("hello gRpc")
	s := NewServer(endpoint)
	s.RegisterService(route.RegisterServer)
	s.Run()
}