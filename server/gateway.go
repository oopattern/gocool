package server

import (
	"google.golang.org/grpc"
	"log"
	"net/http"
	"strings"
)

func gRpcHandlerFunc(s *grpc.Server, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			s.ServeHTTP(w, r)
		} else {
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
