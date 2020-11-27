package main

import (
	"fmt"
	"google.golang.org/grpc"
)

func main() {
	fmt.Print("client connect to server")
	var opts []grpc.DialOption
	conn, err := grpc
}