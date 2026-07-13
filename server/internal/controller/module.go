package controller

import (
	"github.com/go-fuego/fuego"
	"go.uber.org/fx"
)

type Route interface {
	RegisterRoutes(router *fuego.Server)
}

var Module = fx.Module("controller",
	fx.Provide(
		fx.Annotate(NewMediaController, fx.As(new(Route)), fx.ResultTags(`group:"routes"`)),
		fx.Annotate(NewAlbumController, fx.As(new(Route)), fx.ResultTags(`group:"routes"`)),
	),
)
