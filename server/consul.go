package server

import (
	kitlog "github.com/go-kit/kit/log"
	kitsd "github.com/go-kit/kit/sd/consul"
	consul "github.com/hashicorp/consul/api"
	"os"
	"strconv"
	"strings"
	"time"
	"fmt"
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
	registrar := kitsd.NewRegistrar(kitClient, r, kitlog.With(logger, "component", "registrar"))
	registrar.Register()
	ZapLogger.Info(fmt.Sprintf("service[%s] register to consul successfully", service))

	time.Sleep(3*time.Second)

	registrar.Deregister()
	ZapLogger.Info(fmt.Sprintf("service[%s] deregister to consul successfully", service))

	return nil
}
