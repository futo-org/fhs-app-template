package test

import (
	"context"
	"log"

	"fhs/app/internal/config"
	"fhs/app/internal/database"
	"fhs/app/internal/database/migrations"
	"fhs/app/internal/repository"

	"github.com/onsi/ginkgo/v2"
	"github.com/pressly/goose/v3"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/extra/bundebug"
	"go.uber.org/fx"
)

var Config = config.Config{AppEnv: "test"}

// ginkgo.SynchronizedBeforeSuite and ginkgo.SynchronizedAfterSuite are used to run setup and teardown code for the entire test suite
// ensuring that it is only executed once across all parallel test nodes. This should be used for any testcontainers (ie postgres)

// Provisions a unique database for each test run and returns a teardown function to clean up the database after the tests are complete.
// THis example uses SQLite so in-memory is fine for each test run
func setupDatabase() (path string, teardown func()) {
	return ":memory:", func() {}
}

func Configure() (repository.Repositories, func()) {
	// pipe goose migration logs
	goose.SetLogger(log.New(ginkgo.GinkgoWriter, "", log.LstdFlags))
	dbPath, teardownDB := setupDatabase()

	cfg := Config
	cfg.DatabasePath = dbPath

	var repos repository.Repositories
	var db *bun.DB

	app := fx.New(
		fx.NopLogger,
		fx.Supply(&cfg),
		database.Module,
		migrations.Module,
		repository.Module,
		fx.Populate(&repos, &db),
	)

	if err := app.Start(context.Background()); err != nil {
		panic(err)
	}

	// pipe all query logs to ginkgo for output
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithWriter(ginkgo.GinkgoWriter),
		bundebug.WithVerbose(true),
	))

	return repos, func() {
		_ = app.Stop(context.Background())
		teardownDB()
	}
}
