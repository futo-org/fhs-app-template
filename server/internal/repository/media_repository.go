package repository

import (
	"context"
	"database/sql"
	"errors"

	"fhs/app/internal/apperr"
	"fhs/app/internal/database/entity"

	"github.com/uptrace/bun"
)

type MediaRepository interface {
	GetMedia(ctx context.Context, id int64) (*entity.Media, error)
	CreateMedia(ctx context.Context, media *entity.Media) (*entity.Media, error)
	GetByAlbumID(ctx context.Context, albumID int64) ([]*entity.Media, error)
}

type mediaRepository struct {
	db *bun.DB
}

func NewMediaRepository(db *bun.DB) MediaRepository {
	return &mediaRepository{db: db}
}

func (r *mediaRepository) GetMedia(ctx context.Context, id int64) (*entity.Media, error) {
	media := new(entity.Media)

	err := r.db.NewSelect().Model(media).Where("m.id = ?", id).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperr.ErrNotFound
		}
		return nil, err
	}

	return media, nil
}

func (r *mediaRepository) CreateMedia(ctx context.Context, media *entity.Media) (*entity.Media, error) {
	newMedia := &entity.Media{
		Title: media.Title,
		URL:   media.URL,
		Path:  media.Path,
	}

	_, err := r.db.NewInsert().Model(newMedia).Returning("*").Exec(ctx)
	if err != nil {
		return nil, err
	}

	return newMedia, nil
}

// GetByAlbumID returns the media linked to an album, joined through the
// album_media link table.
func (r *mediaRepository) GetByAlbumID(ctx context.Context, albumID int64) ([]*entity.Media, error) {
	var media []*entity.Media

	err := r.db.NewSelect().
		Model(&media).
		Join("JOIN album_media AS am ON am.media_id = m.id").
		Where("am.album_id = ?", albumID).
		Order("m.id ASC").
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return media, nil
}
