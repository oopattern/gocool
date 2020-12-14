package server

import (
	"context"
	"fmt"
	"log"
	"strings"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/oopattern/gocool/proto"
	"google.golang.org/grpc"
	"net/http"
)

func grpcHandlerFunc(s *grpc.Server, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ZapLogger.Info(fmt.Sprintf("coming request[%+v]", *r))
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			ZapLogger.Debug("zzzzz?")
			s.ServeHTTP(w, r)
		} else {
			ZapLogger.Debug("wwww?")
			h.ServeHTTP(w, r)
		}
	})
}

// support REST gateway
func RouteHttp(gwAddr string, endpoint string, s *grpc.Server) error {
	log.Printf("addr[%s] endpoint[%s] route to gateway", gwAddr, endpoint)
	gwMux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := proto.RegisterObserveHandlerFromEndpoint(context.Background(), gwMux, endpoint, opts); err != nil {
		ZapLogger.Error("register gateway err")
		return err
	}
	mux := http.NewServeMux()
	mux.Handle("/", gwMux)
	h := &http.Server{
		Addr:		gwAddr,
		Handler:	grpcHandlerFunc(s, mux),
	}

	go h.ListenAndServe()

	ZapLogger.Info("register to gateway successfully")
	return nil
}
