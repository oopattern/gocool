package main

import (
	"context"
	"fmt"
	"log"
	"google.golang.org/grpc"
	"github.com/oopattern/gocool/proto"
)

var (
	RouteHost = "localhost"
	RoutePort = 7777
)

func main() {
	addr := fmt.Sprintf("%s:%d", RouteHost, RoutePort)
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()
	client := proto.NewObserveClient(conn)
	req := proto.RouteReq{
		Name: "Sakulali",
	}
	resp, err := client.SayRoute(context.Background(), &req)
	if nil == err {
		fmt.Printf("client get resp ip[%s] port[%s] ok\n", resp.GetIp(), resp.GetPort())
	}
}