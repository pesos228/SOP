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

func Seed(cfg database.Config, timeOut time.Duration) error {
	db, err := sql.Open("pgx", cfg.DSN())
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	fmt.Println("Seeding data...")

	if err := kitMigrate.Migrate(ctx, db, migration.EmbedSeeds, migration.SeedsDir); err != nil {
		return fmt.Errorf("seed database: %w", err)
	}

	fmt.Println("Seed data completed!")

	return nil
}
