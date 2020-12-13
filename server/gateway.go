package server

import (
	"fmt"
	"log"
	"context"
	"strings"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/oopattern/gocool/proto"
	"google.golang.org/grpc"
	"net/http"
)

func RouteHttp(endpoint string, s *grpc.Server) *http.Server {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := proto.RegisterObserveHandlerFromEndpoint(context.Background(), mux, endpoint, opts); err != nil {
		log.Fatalf("register handler err[%v]", err)
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ZapLogger.Info(fmt.Sprintf("request[%+v] coming", *r))
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			s.ServeHTTP(w, r)
		}
	})
	ZapLogger.Info("register to gateway successfully")
	return &http.Server{
		Addr:              endpoint,
		Handler:           handler,
	}
}
