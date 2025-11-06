package app

import (
	"fmt"
	"hosting-service/internal/domain"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func (a *App) initDB() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(a.config.DB_DSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connection to DB failed: %w", err)
	}

	log.Println("Running DB migrations...")
	db.AutoMigrate(&domain.Plan{}, &domain.Server{})

	return db, nil
}
