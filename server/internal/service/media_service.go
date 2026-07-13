package service

import (
	"context"
	"fhs/app/internal/database/entity"
	"fhs/app/internal/dto"
	"fhs/app/internal/repository"
)

type MediaService interface {
	GetMedia(ctx context.Context, id int64) (*dto.MediaDto, error)
	CreateMedia(ctx context.Context, dto *dto.CreateMediaDto) (*dto.MediaDto, error)
}

type mediaService struct {
	repos repository.Repositories
}

func NewMediaService(repos repository.Repositories) *mediaService {
	return &mediaService{repos: repos}
}

func (s *mediaService) GetMedia(ctx context.Context, id int64) (*dto.MediaDto, error) {
	media, err := s.repos.MediaRepository.GetMedia(ctx, id)
	if err != nil {
		return nil, err
	}

	return mapMedia(media), nil
}

func (s *mediaService) CreateMedia(ctx context.Context, dto *dto.CreateMediaDto) (*dto.MediaDto, error) {
	createdMedia, err := s.repos.MediaRepository.CreateMedia(ctx, &entity.Media{
		Title: dto.Name,
		URL:   dto.Path,
		Path:  dto.Path,
	})

	if err != nil {
		return nil, err
	}

	return mapMedia(createdMedia), nil
}

func mapMedia(media *entity.Media) *dto.MediaDto {
	return &dto.MediaDto{
		ID:    media.ID,
		Title: media.Title,
		URL:   media.URL,
		Path:  media.Path,
	}
}
