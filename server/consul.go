package server

import (
	"errors"
	"github.com/hashicorp/consul/agent/consul"
	kitlog "github.com/go-kit/kit/log"
	kitsd "github.com/go-kit/kit/sd/consul"
	"github.com/oopattern/gocool/config"
	"github.com/oopattern/gocool/log"
	"os"
	"strconv"
	"strings"
)

// # nohup consul agent -dev -client=0.0.0.0 &

var (
	consulRegistrar *kitsd.Registrar
)

func DeregisterConsul(service string) error {
	if consulRegistrar == nil {
		return errors.New("consul registrar is nil error")
	}
	consulRegistrar.Deregister()
	log.Info("service[%s] deregister to consul successfully", service)
	return nil
}

func RegisterConsul(service string, endpoint string) error {
	addr := strings.Split(endpoint, ":")
	port, _ := strconv.Atoi(addr[1])

	client, err := consul.NewClient(&consul.Config{Address:config.ConsulEndPoint})
	if err != nil {
		log.Error("create consul client error")
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
	consulRegistrar = kitsd.NewRegistrar(kitClient, r, kitlog.With(logger, "component", "registrar"))
	consulRegistrar.Register()

	log.Info("service[%s] register to consul successfully", service)
	return nil
}
