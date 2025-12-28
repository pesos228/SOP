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

func Migrate(cfg database.Config, timeOut time.Duration) error {
	db, err := sql.Open("pgx", cfg.DSN())
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	fmt.Println("Applying migrations...")

	if err := kitMigrate.Migrate(ctx, db, migration.EmbedMigrations, migration.MigrationsDir); err != nil {
		return fmt.Errorf("migrate up failed: %w", err)
	}

	fmt.Println("Migrations applied successfully!")

	return nil
}
