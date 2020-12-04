package server

import (
	"fmt"
	"io"
	"os"
	"github.com/go-kit/kit/sd"
	kitlog "github.com/go-kit/kit/log"
	kitsd "github.com/go-kit/kit/sd/consul"
	kitep "github.com/go-kit/kit/endpoint"
	consul "github.com/hashicorp/consul/api"
	"strconv"
	"strings"
)

var (
	// # nohup consul agent -dev -client=0.0.0.0 &
	consulAddr = "localhost:8500"
)

func RegisterConsul(service string, endpoint string) error {
	addr := strings.Split(endpoint, ":")
	port, _ := strconv.Atoi(addr[1])

	client, err := consul.NewClient(&consul.Config{Address:consulAddr})
	if err != nil {
		ZapLogger.Error("create consul client error")
		return err
	}
	kitClient := kitsd.NewClient(client)
	logger := kitlog.NewLogfmtLogger(os.Stderr)

	r := &consul.AgentServiceRegistration{
		ID:                "Kit-Consul",
		Name:              service,
		Tags:              []string{"alpha"},
		Port:              port,
		Address:           addr[0],
		EnableTagOverride: false,
	}
	factory := func(instance string) (kitep.Endpoint, io.Closer, error) {
		return kitep.Nop, nil, nil
	}
	instancer := kitsd.NewInstancer(
		kitClient,
		kitlog.With(logger, "component", "instancer"),
		r.Name,
		r.Tags,
		true,
	)
	var _ = sd.NewEndpointer(instancer, factory, kitlog.With(logger, "component", "endpointer"))
	registrar := kitsd.NewRegistrar(kitClient, r, kitlog.With(logger, "component", "registrar"))
	registrar.Register()

	ZapLogger.Info(fmt.Sprintf("service[%s] register to consul successfully", service))
	return nil
}
