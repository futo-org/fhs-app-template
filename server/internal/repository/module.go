package repository

import "go.uber.org/fx"

type Repositories struct {
	fx.In
	MediaRepository MediaRepository
	AlbumRepository AlbumRepository
}

var Module = fx.Module("repository", fx.Provide(
	NewMediaRepository,
	NewAlbumRepository,
))
