package main

import (
	"fmt"
	"log"

	"gitlab.k8s.gromnsk.ru/workshop/austin/pkg/config"
	"gitlab.k8s.gromnsk.ru/workshop/austin/pkg/tracking"
)

func main() {
	fmt.Println("Tracking workshop")
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("filet to get config: %v\n", err)
	}
	fmt.Printf("%+v", cfg)

	tracking.Run(cfg)
}
