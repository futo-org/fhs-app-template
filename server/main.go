package main

import (
	"fhs/app/internal/config"
	"fhs/app/internal/controller"
	"fhs/app/internal/database"
	"fhs/app/internal/database/migrations"
	"fhs/app/internal/repository"
	"fhs/app/internal/server"
	"fhs/app/internal/service"

	"go.uber.org/fx"
)

// Main
func main() {
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
