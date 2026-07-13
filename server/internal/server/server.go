package server

import (
	"context"
	"errors"
	"fhs/app/internal/config"
	"fhs/app/internal/controller"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-fuego/fuego"
	"go.uber.org/fx"
)

type RouteParam struct {
	fx.In
	Routes []controller.Route `group:"routes"`
}

// All API routes live under /api so they can never collide with the static
// web UI mounted on "/" (see webui.go)
func RegisterRoutes(router *fuego.Server, rp RouteParam) {
	api := fuego.Group(router, "/api")
	for _, route := range rp.Routes {
		route.RegisterRoutes(api)
	}
}

func NewServer(cfg *config.Config) *fuego.Server {
	isDev := cfg.AppEnv == "development"

	return fuego.NewServer(
		fuego.WithAddr(cfg.ServerAddr),
		fuego.WithEngineOptions(
			fuego.WithOpenAPIConfig(fuego.OpenAPIConfig{
				JSONFilePath:     "doc/openapi.json",
				PrettyFormatJSON: true,
				DisableLocalSave: !isDev,
				DisableMessages:  !isDev,
				DisableSwaggerUI: !isDev,
				Info: &openapi3.Info{
					Title:   "FHS Server Template",
					Version: "0.1.0",
				},
			}),
		),
	)
}

func RunServer(lc fx.Lifecycle, s *fuego.Server) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := s.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					panic(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return s.Shutdown(ctx)
		},
	})
}
