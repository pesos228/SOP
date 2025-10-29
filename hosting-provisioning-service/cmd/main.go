package main

import (
	"hosting-provisioning-service/internal/config"
	"hosting-provisioning-service/internal/worker"
	"log"
)

func main() {
	cfg := config.Load()

	appWorker, err := worker.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize worker: %v", err)
	}

	appWorker.Run()
}
