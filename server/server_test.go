package server

import (
	"fmt"
	"testing"
)

var (
	route = RouteServer{}
	endpoint = fmt.Sprintf("localhost:%d",7777)
)

func TestBuildGrpcServer(t *testing.T) {
	ZapLogger.Info(fmt.Sprintf("trace port[%d]", tracePort))
	ZapLogger.Debug("hello gRpc")
	s := NewServer(endpoint)
	s.RegisterService(route.RegisterServer)
	s.Run()
}