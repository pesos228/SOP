package main

import (
	"hosting-provisioning-service/internal/app"
	"hosting-provisioning-service/internal/config"
)

func main() {
	cfg := config.Load()

	app := app.New(cfg)

	app.Run()
}
