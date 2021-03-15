package server

import (
	"errors"
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

// 查找服务
func LookupService(dc, name string) ([]*consul.ServiceEntry, error) {
	client, err := consul.NewClient(&consul.Config{Address:config.ConsulEndPoint})
	if err != nil {
		log.Error("create consul client error")
		return nil, err
	}

	opt := &consul.QueryOptions{
		Datacenter:        dc,
		Filter:            "",
	}
	kitClient := kitsd.NewClient(client)
	entrys, _, err := kitClient.Service(name, "", false, opt)
	if err != nil {
		log.Error("consul get service info err[%+v]", err)
		return nil, err
	}
	if len(entrys) <= 0 {
		log.Info("consul query service[%s] is empty", name)
		return nil, errors.New("service is not found error")
	}
	for _, item := range entrys {
		log.Info("show service agent[%+v]", *item.Service)
	}

	return entrys, nil
}

// 注销consul的服务发现
func (s *grpcServer) deregisterConsul() {
	for service, registrar := range consulRegistrars {
		registrar.Deregister()
		log.Info("service_name[%s] deregister to consul successfully", service)
	}
}

// 注册服务到consul
func (s *grpcServer) registerConsul() error {
	// register to consul
	endpoint := s.GetListener().Addr().String()
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

		// 如果一个服务需要注册多个实例ip-port, 需要注意AgentServiceRegistration的ID参数不能相同
		r := &consul.AgentServiceRegistration{
			ID:                service + "-" + endpoint,
			Name:              service,
			Tags:              []string{"alpha"},
			Port:              port,
			Address:           addr[0],
			EnableTagOverride: false,
		}
		registrar := kitsd.NewRegistrar(kitClient, r, kitlog.With(logger, "component", "registrar"))
		registrar.Register()
		consulRegistrars[service] = registrar

		// 查询已注册的服务
		var _, _ = LookupService(config.DefaultConsulDataCenter, service)

		log.Info("service_name[%s] endpoint[%s] info[%+v] register to consul successfully", service, endpoint, info)
	}

	return nil
}
