package migrations

import (
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"
	"go.uber.org/fx"
)

//go:embed *.sql
var migrationsFS embed.FS

func runMigration(db *sql.DB) error {

	goose.SetBaseFS(migrationsFS)
	goose.SetDialect("sqlite3")

	return goose.Up(db, ".")
}

var Module = fx.Module("migrations", fx.Invoke(runMigration))
