package server

import (
	"fmt"
	"log"
	"context"
	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/reporter"
)

var (
	// skywalking backend port
	TracePort = 11800
)

func init() {
	endpoint := fmt.Sprintf("127.0.0.1:%d", TracePort)
	r, err := reporter.NewGRPCReporter(endpoint)
	if err != nil {
		log.Fatalf("failed new trace reporter err: %v", err)
	}
	tracer, err := go2sky.NewTracer("SayRoute", go2sky.WithReporter(r))
	if err != nil {
		log.Fatalf("failed create tracer err: %v", err)
	}
	span, _, err := tracer.CreateLocalSpan(context.Background())
	if err != nil {
		log.Fatalf("failed create local span err: %v", err)
	}
	span.End()
}

