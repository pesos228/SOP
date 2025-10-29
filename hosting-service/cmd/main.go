package main

import (
	"hosting-service/internal/app"
	"hosting-service/internal/config"
)

func main() {
	cfg := config.Load()

	application := app.New(cfg)

	application.Run()
}
