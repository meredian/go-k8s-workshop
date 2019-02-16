package discovery

import (
	"errors"
	"fmt"

	"github.com/hashicorp/consul/api"
	"gitlab.k8s.gromnsk.ru/workshop/austin/pkg/config"
)

type ServiceNode struct {
	ID      string
	Address string
	Port    int
}

func NewConsulClient(hostname string) (*api.Client, error) {
	consulConfig := &api.Config{}
	consulConfig.Address = hostname

	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		return nil, err
	}

	return consulClient, nil
}

func GetServiceNode(consulClient *api.Client, serviceName string) (*ServiceNode, error) {
	serviceEntry, _, err := consulClient.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, err
	}

	if len(serviceEntry) == 0 {
		return nil, errors.New(fmt.Sprintf("service %s is not discovered\n", serviceName))
	}

	sn := &ServiceNode{
		ID:      serviceEntry[0].Service.ID,
		Address: serviceEntry[0].Service.Address,
		Port:    serviceEntry[0].Service.Port,
	}

	return sn, nil
}

func RegisterService(consulClient *api.Client, cfg *config.Config) error {
	err := consulClient.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      getServiceID(cfg.Consul.Servicename, cfg.Server.Host),
		Name:    cfg.Consul.Servicename,
		Address: cfg.Server.Host,
		Port:    cfg.Server.Port,
	})

	if err != nil {
		return err
	}

	return nil
}

func DeregisterService(consulClient *api.Client, cfg *config.Config) error {
	err := consulClient.Agent().ServiceDeregister(getServiceID(cfg.Consul.Servicename, cfg.Server.Host))
	if err != nil {
		return err
	}

	return nil
}

func getServiceID(serviceName, hostname string) string {
	return fmt.Sprintf("%s-%s", serviceName, hostname)
}
