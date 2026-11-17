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
	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	return goose.Up(db, ".")
}

var Module = fx.Module("migrations", fx.Invoke(runMigration))
