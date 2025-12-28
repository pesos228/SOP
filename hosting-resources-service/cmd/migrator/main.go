package main

import (
	"errors"
	"fmt"
	"hosting-kit/database"
	"hosting-resources-service/cmd/migrator/commands"
	"log"
	"os"
	"time"

	"github.com/ardanlabs/conf/v3"
)

func main() {
	if err := run(); err != nil {
		log.Printf("error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg := struct {
		Args conf.Args
		DB   struct {
			User     string `conf:"default:postgres"`
			Password string `conf:"default:vladick,mask"`
			Host     string `conf:"default:localhost:5432"`
			Name     string `conf:"default:sop_pool"`
		}
		Migration struct {
			Timeout time.Duration `conf:"default:10s"`
		}
	}{}

	const prefix = "RES"

	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			os.Exit(0)
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	dbConfig := database.Config{
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		Host:     cfg.DB.Host,
		Name:     cfg.DB.Name,
	}

	return processCommands(cfg.Args, cfg.Migration.Timeout, dbConfig)
}

func processCommands(args conf.Args, timeOut time.Duration, dbConfig database.Config) error {
	switch args.Num(0) {
	case "migrate", "up":
		return commands.Migrate(dbConfig, timeOut)

	case "rollback", "down":
		return commands.Rollback(dbConfig, timeOut)

	case "seed":
		return commands.Seed(dbConfig, timeOut)

	case "migrate-seed", "up-seed":
		if err := commands.Migrate(dbConfig, timeOut); err != nil {
			return err
		}
		if err := commands.Seed(dbConfig, timeOut); err != nil {
			return err
		}
		return nil

	case "status":
		return commands.Status(dbConfig, timeOut)

	case "reset":
		return commands.Reset(dbConfig, timeOut)

	default:
		fmt.Println("migrate/up:         create the schema in the database")
		fmt.Println("rollback/down:      roll back the most recent migration")
		fmt.Println("seed:               seed the database with initial data")
		fmt.Println("migrate-seed/up-seed: run migrations then seed")
		fmt.Println("status:             print the status of all migrations")
		fmt.Println("reset:              roll back all migrations")

		return errors.New("unknown command")
	}
}
