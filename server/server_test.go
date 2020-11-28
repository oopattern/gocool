package server

import (
	"fmt"
	"testing"
)

var (
	endpoint = fmt.Sprintf("%s:%d", "localhost", 7777)
)

func TestBuildGrpcServer(t *testing.T) {
	fmt.Println("hello gRpc")
	route := RouteServer{}
	s := NewServer(endpoint)
	s.RegisterService(route.RegisterServer)
	s.Run()
}