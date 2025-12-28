package migration

import (
	"embed"
)

//go:embed sql/*.sql
var EmbedMigrations embed.FS

//go:embed seeds/*.sql
var EmbedSeeds embed.FS

const (
	MigrationsDir = "sql"
	SeedsDir      = "seeds"
)
