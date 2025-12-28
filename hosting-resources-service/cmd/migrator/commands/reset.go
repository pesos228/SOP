package commands

import (
	"context"
	"database/sql"
	"fmt"
	"hosting-kit/database"
	kitMigrate "hosting-kit/migration"
	"hosting-resources-service/internal/platform/migration"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Reset(cfg database.Config, timeOut time.Duration) error {
	db, err := sql.Open("pgx", cfg.DSN())
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	fmt.Println("Resetting all migrations...")

	if err := kitMigrate.Reset(ctx, db, migration.EmbedMigrations, migration.MigrationsDir); err != nil {
		return fmt.Errorf("migrate reset: %w", err)
	}

	fmt.Println("Reset successful!")

	return nil
}
