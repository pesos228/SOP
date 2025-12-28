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

func Status(cfg database.Config, timeOut time.Duration) error {
	db, err := sql.Open("pgx", cfg.DSN())
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	if err := kitMigrate.Status(ctx, db, migration.EmbedMigrations, migration.MigrationsDir); err != nil {
		return fmt.Errorf("migrate status failed: %w", err)
	}

	return nil
}
