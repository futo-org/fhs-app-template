package database

import (
	"context"
	"database/sql"
	"fhs/app/internal/config"
	"fhs/app/internal/database/entity"
	"os"
	"path/filepath"

	_ "github.com/glebarez/go-sqlite"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/extra/bundebug"
	"go.uber.org/fx"
)

func SetupDatabase(lc fx.Lifecycle, cfg *config.Config) (*sql.DB, error) {
	if dir := filepath.Dir(cfg.DatabasePath); dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}
	}

	db, err := sql.Open("sqlite", cfg.DatabasePath)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return db.PingContext(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return db.Close()
		},
	})

	return db, nil
}

func NewBun(db *sql.DB, cfg *config.Config) *bun.DB {
	bunDB := bun.NewDB(db, sqlitedialect.New())

	bunDB.RegisterModel((*entity.AlbumMedia)(nil))

	if cfg.AppEnv == "development" {
		bunDB.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	return bunDB
}

var Module = fx.Module("database", fx.Provide(SetupDatabase, NewBun))
