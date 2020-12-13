package server

import (
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"net/http"
)

func RouteHttp(endpoint string, s *grpc.Server) *http.Server {
	opts := []grpc.DialOption{grpc.WithInsecure()}
	mux := runtime.NewServeMux()
	return &http.Server{
		Addr:              endpoint,
		Handler:           s.ServeHTTP,
	}
}
