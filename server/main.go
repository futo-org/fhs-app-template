package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
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

	generateOpenAPI := flag.Bool("open-api", false, "print the OpenAPI spec and exit")
	flag.Parse()

	modules := fx.Options(
		database.Module,
		repository.Module,
		service.Module,
		controller.Module,
		server.Module,
	)

	if *generateOpenAPI {
		if err := exportOpenAPI(modules); err != nil {
			log.Fatalf("export OpenAPI: %v", err)
		}

		return
	}

	fmt.Printf("version is %s\n", version.Version)

	fx.New(config.Module, migrations.Module, modules).Run()
}
