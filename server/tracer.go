package server

import (
	"context"
	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/reporter"
	"github.com/oopattern/gocool/config"
	"github.com/oopattern/gocool/log"
)

func init() {
	r, err := reporter.NewGRPCReporter(config.SkyWalkingEndPoint)
	if err != nil {
		log.Fatal("failed new trace reporter err: %+v", err)
	}
	tracer, err := go2sky.NewTracer("SayRoute", go2sky.WithReporter(r))
	if err != nil {
		log.Fatal("failed create tracer err: %+v", err)
	}
	span, _, err := tracer.CreateLocalSpan(context.Background())
	if err != nil {
		log.Fatal("failed create local span err: %+v", err)
	}
	span.End()
}

