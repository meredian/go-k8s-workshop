package tracking

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.k8s.gromnsk.ru/workshop/austin/pkg/discovery"

	"gitlab.k8s.gromnsk.ru/workshop/austin/pkg/config"
	"gitlab.k8s.gromnsk.ru/workshop/austin/pkg/handlers"
	"gitlab.k8s.gromnsk.ru/workshop/austin/pkg/storage"
)

func Run(cfg *config.Config) {

	consulClient, err := discovery.NewConsulClient(cfg.Consul.Hostport)
	if err != nil {
		log.Fatalf("failed to connect consul client: %v\n", err)
	}

	serviceNode, err := discovery.GetServiceNode(consulClient, "cassandra")
	if err != nil {
		log.Fatalf("failed to get service node: %v\n", err)
	}

	// FIXME: HOST should be serviceNode.Address
	hostport := fmt.Sprintf("%s:%d", "127.0.0.1", serviceNode.Port)
	keyspace := cfg.Storage.Keyspace
	session, err := storage.InitStorage(hostport, keyspace)
	if err != nil {
		log.Fatalf("failed to init storage: %v\n", err)
	}

	server := handlers.Server{
		Session: session,
	}

	err = discovery.RegisterService(consulClient, cfg)
	if err != nil {
		log.Fatalf("failed to register service: %v\n", err)
	}
	defer func() {
		err := discovery.DeregisterService(consulClient, cfg)
		if err != nil {
			log.Printf("failed to deregister service: %v\n", err)
		}
	}()

	handlers.RegisterMetrics()

	http.HandleFunc("/save", server.SaveActionHandler)
	http.HandleFunc("/get-status", server.GetActionStatusHandler)

	http.HandleFunc("/healthz", server.HealthHandler)
	http.Handle("/metrics", promhttp.Handler())

	fmt.Printf("Service starting on port %d..\n", cfg.Server.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), nil)
	if err != nil {
		log.Fatalf("failed to start server: %v\n", err)
	}
}
