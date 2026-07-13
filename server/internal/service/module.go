package service

import (
	"go.uber.org/fx"
)

var Module = fx.Module("service",
	fx.Provide(
		fx.Annotate(NewMediaService, fx.As(new(MediaService))),
		fx.Annotate(NewAlbumService, fx.As(new(AlbumService))),
	),
)
