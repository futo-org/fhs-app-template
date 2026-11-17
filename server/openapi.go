package main

import (
	"context"
	"encoding/json"
	"os"

	"fhs/app/internal/config"

	"github.com/go-fuego/fuego"
	"go.uber.org/fx"
)

func exportOpenAPI(modules fx.Option) error {
	var router *fuego.Server
	app := fx.New(
		modules,
		fx.NopLogger,
		fx.Supply(&config.Config{DatabasePath: ":memory:"}),
		fx.Populate(&router),
	)
	if err := app.Err(); err != nil {
		return err
	}

	spec := router.OutputOpenAPISpec()
	if err := spec.Validate(context.Background()); err != nil {
		return err
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "\t")
	return encoder.Encode(spec)
}
