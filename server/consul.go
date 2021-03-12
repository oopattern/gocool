package server

import (
	"os"
	"strconv"
	"strings"
	"time"
	kitlog "github.com/go-kit/kit/log"
	kitsd "github.com/go-kit/kit/sd/consul"
	consul "github.com/hashicorp/consul/api"
	"github.com/oopattern/gocool/log"
	"github.com/oopattern/gocool/config"
)

// # nohup consul agent -dev -client=0.0.0.0 &

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
	registrar := kitsd.NewRegistrar(kitClient, r, kitlog.With(logger, "component", "registrar"))
	registrar.Register()
	log.Info("service[%s] register to consul successfully", service)

	time.Sleep(3*time.Second)

	registrar.Deregister()
	log.Info("service[%s] deregister to consul successfully", service)

	return nil
}
