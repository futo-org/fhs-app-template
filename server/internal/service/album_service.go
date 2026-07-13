package service

import (
	"context"

	"fhs/app/internal/config"
	"fhs/app/internal/database/entity"
	"fhs/app/internal/dto"
	"fhs/app/internal/repository"
)

type AlbumService interface {
	CreateAlbum(ctx context.Context, in *dto.CreateAlbumDto) (*dto.AlbumDto, error)
	GetAlbum(ctx context.Context, id int64) (*dto.AlbumDto, error)
	ListAlbums(ctx context.Context) ([]*dto.AlbumDto, error)
	UpdateAlbum(ctx context.Context, id int64, in *dto.UpdateAlbumDto) (*dto.AlbumDto, error)
	DeleteAlbum(ctx context.Context, id int64) error
	AddMedia(ctx context.Context, albumID, mediaID int64) error
	RemoveMedia(ctx context.Context, albumID, mediaID int64) error
}

type albumService struct {
	repos  repository.Repositories
	config *config.Config
}

func NewAlbumService(repos repository.Repositories, config *config.Config) AlbumService {
	return &albumService{repos: repos, config: config}
}

func (s *albumService) CreateAlbum(ctx context.Context, in *dto.CreateAlbumDto) (*dto.AlbumDto, error) {
	album, err := s.repos.AlbumRepository.Create(ctx, &entity.Album{Name: in.Name})
	if err != nil {
		return nil, err
	}
	return mapAlbum(album, nil), nil
}

func (s *albumService) GetAlbum(ctx context.Context, id int64) (*dto.AlbumDto, error) {
	album, err := s.repos.AlbumRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Fetch the album's media through the link table (media repo owns this query).
	media, err := s.repos.MediaRepository.GetByAlbumID(ctx, id)
	if err != nil {
		return nil, err
	}

	return mapAlbum(album, media), nil
}

func (s *albumService) ListAlbums(ctx context.Context) ([]*dto.AlbumDto, error) {
	albums, err := s.repos.AlbumRepository.List(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.AlbumDto, 0, len(albums))
	for _, album := range albums {
		result = append(result, mapAlbum(album, nil)) // list is shallow (no media) to avoid N+1
	}
	return result, nil
}

func (s *albumService) UpdateAlbum(ctx context.Context, id int64, in *dto.UpdateAlbumDto) (*dto.AlbumDto, error) {
	if err := s.repos.AlbumRepository.Update(ctx, &entity.Album{ID: id, Name: in.Name}); err != nil {
		return nil, err
	}
	return s.GetAlbum(ctx, id)
}

func (s *albumService) DeleteAlbum(ctx context.Context, id int64) error {
	return s.repos.AlbumRepository.Delete(ctx, id)
}

func (s *albumService) AddMedia(ctx context.Context, albumID, mediaID int64) error {
	return s.repos.AlbumRepository.AddMedia(ctx, albumID, mediaID)
}

func (s *albumService) RemoveMedia(ctx context.Context, albumID, mediaID int64) error {
	return s.repos.AlbumRepository.RemoveMedia(ctx, albumID, mediaID)
}

func mapAlbum(album *entity.Album, media []*entity.Media) *dto.AlbumDto {
	mediaDtos := make([]dto.MediaDto, 0, len(media))
	for _, m := range media {
		mediaDtos = append(mediaDtos, *mapMedia(m))
	}

	return &dto.AlbumDto{
		ID:    album.ID,
		Name:  album.Name,
		Media: mediaDtos,
	}
}
