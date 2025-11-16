package app

import (
	"fmt"

	"hosting-service/internal/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func (a *App) initDB() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(a.config.DB_DSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connection to DB failed: %w", err)
	}

	err = db.AutoMigrate(&domain.Plan{}, &domain.Server{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}
