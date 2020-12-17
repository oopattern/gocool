package server

import (
	"fmt"
	"log"
	"google.golang.org/grpc"
	"net/http"
	"strings"
)

func gRpcHandlerFunc(s *grpc.Server, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ZapLogger.Info(fmt.Sprintf("coming request[%+v]", *r))
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			ZapLogger.Debug("zzzzz?")
			s.ServeHTTP(w, r)
		} else {
			ZapLogger.Debug("wwwww?")
			h.ServeHTTP(w, r)
		}
	})
}

// support REST gateway
func StartGateway(gwAddr string, s *grpc.Server) {
	mux := http.NewServeMux()
	mux.Handle("/", GatewayMux)
	h := &http.Server{
		Addr:		gwAddr,
		Handler:	gRpcHandlerFunc(s, mux),
	}
	go func() {
		if err := h.ListenAndServe(); err != nil {
			log.Fatalf("failed to lister http")
		}
	}()
}
