package commands

import (
	"context"
	"database/sql"
	"fmt"
	"hosting-kit/database"
	kitMigrate "hosting-kit/migration"
	"hosting-service/internal/platform/migration"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Rollback(cfg database.Config, timeOut time.Duration) error {
	db, err := sql.Open("pgx", cfg.DSN())
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	fmt.Println("Rolling back last migration...")

	if err := kitMigrate.Rollback(ctx, db, migration.EmbedMigrations, migration.MigrationsDir); err != nil {
		return fmt.Errorf("migrate down failed: %w", err)
	}

	fmt.Println("Rollback successful!")

	return nil
}
