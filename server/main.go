package main

import (
	_ "embed"
	"fmt"
	"os"
	"strings"

	"fhs/app/internal/config"
	"fhs/app/internal/controller"
	"fhs/app/internal/database"
	"fhs/app/internal/database/migrations"
	"fhs/app/internal/repository"
	"fhs/app/internal/server"
	"fhs/app/internal/service"
	"fhs/app/internal/version"

	"go.uber.org/fx"
)

// The VERSION file is the single source of truth for the app version; it is
// embedded so every build (dev, test, prod) reports the same version
//
//go:embed VERSION
var rawVersion string

// Main
func main() {
	version.Version = strings.TrimSpace(rawVersion)

	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Println(version.Version)
		return
	}

	fmt.Printf("version is %s\n", version.Version)

	fx.New(
		config.Module,
		database.Module,
		migrations.Module,
		repository.Module,
		service.Module,
		controller.Module,
		server.Module,
	).Run()
}
