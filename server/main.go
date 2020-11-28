package main

import (
	"fmt"
	"github.com/oopattern/gocool/server"
)

var (
	endpoint = fmt.Sprintf("%s:%d", "localhost", 7777)
)

func main() {
	fmt.Println("hello gRpc")
	route := server.RouteServer{}
	s := server.NewServer(endpoint)
	s.RegisterService(route.RegisterServer)
	s.Run()
}
