package tracking

import (
	"fmt"
	"log"
	"net/http"

	"gitlab.k8s.gromnsk.ru/workshop/austin/pkg/config"
	"gitlab.k8s.gromnsk.ru/workshop/austin/pkg/handlers"
	"gitlab.k8s.gromnsk.ru/workshop/austin/pkg/storage"
)

func Run(cfg *config.Config) {

	hostport := fmt.Sprintf("%s:%d", cfg.Storage.Host, cfg.Storage.Port)
	keyspace := cfg.Storage.Keyspace
	session, err := storage.InitStorage(hostport, keyspace)
	if err != nil {
		log.Fatalf("failed to init storage: %v\n", err)
	}

	server := handlers.Server{
		Session: session,
	}

	http.HandleFunc("/save", server.SaveActionHandler)
	http.HandleFunc("/get-status", server.GetActionStatusHandler)

	http.HandleFunc("/healthz", server.HealthHandler)

	fmt.Printf("Service starting on port %d..\n", cfg.Server.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), nil)
	if err != nil {
		log.Fatalf("failed to start server: %v\n", err)
	}
}
