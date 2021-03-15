package server

import (
	kitlog "github.com/go-kit/kit/log"
	kitsd "github.com/go-kit/kit/sd/consul"
	consul "github.com/hashicorp/consul/api"
	"github.com/oopattern/gocool/config"
	"github.com/oopattern/gocool/log"
	"os"
	"strconv"
	"strings"
)

// # nohup consul agent -dev -client=0.0.0.0 &

var (
	// service_name --> registrar
	consulRegistrars = make(map[string]*kitsd.Registrar)
)

func (s *grpcServer) deregisterConsul() {
	for service, registrar := range consulRegistrars {
		registrar.Deregister()
		log.Info("service_name[%s] deregister to consul successfully", service)
	}
}

func (s *grpcServer) registerConsul() error {
	// register to consul
	endpoint := s.listener.Addr().String()
	for service, info := range s.server.GetServiceInfo() {
		// 服务已注册
		if _, ok := consulRegistrars[service]; ok {
			continue
		}

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
			ID:                service,
			Name:              service,
			Tags:              []string{"alpha"},
			Port:              port,
			Address:           addr[0],
			EnableTagOverride: false,
		}
		registrar := kitsd.NewRegistrar(kitClient, r, kitlog.With(logger, "component", "registrar"))
		registrar.Register()
		consulRegistrars[service] = registrar

		log.Info("service_name[%s] endpoint[%s] info[%+v] register to consul successfully", service, endpoint, info)
	}

	return nil
}
