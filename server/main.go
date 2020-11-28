package main

import (
	"fmt"
)

var (
	RouteHost = "localhost"
	RoutePort = 7777
)

func main() {
	fmt.Println("hello gRpc")
	endpoint := fmt.Sprintf("%s:%d", RouteHost, RoutePort)
	grpcServer := NewServer()
	routeServer := routeServer{}
	routeServer.RegisterServer(grpcServer)
	Run(endpoint, grpcServer)
}
