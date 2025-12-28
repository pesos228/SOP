package migration

import (
	"embed"
)

//go:embed sql/*.sql
var EmbedMigrations embed.FS

const MigrationsDir = "sql"
