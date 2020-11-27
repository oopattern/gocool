package gocool

import (
	"fmt"
	"context"
	"gocool/proto"
)

type routeServer struct {}

// rpc SayRoute(RouteReq) returns (RouteResp) {}
func (r *routeServer) SayRoute(ctx context.Context, req *proto.RouteReq) (*proto.RouteResp, error) {
	fmt.Println(req.GetName())
	return &proto.RouteResp{
		Ip:   "localhost",
		Port: "7777",
	}, nil
}

func main() {
	fmt.Println("hello world")
	r := routeServer{}
	req := proto.RouteReq{
		Name: "RouteName",
	}
	resp, err := r.SayRoute(context.Background(), &req)
	if err == nil {
		fmt.Printf("ip[%s] port[%s]\n", resp.GetIp(), resp.GetPort())
	}
}
