package repository

import (
	"context"
	"database/sql"
	"errors"

	"fhs/app/internal/apperr"
	"fhs/app/internal/database/entity"

	"github.com/uptrace/bun"
)

type AlbumRepository interface {
	Create(ctx context.Context, album *entity.Album) (*entity.Album, error)
	GetByID(ctx context.Context, id int64) (*entity.Album, error)
	List(ctx context.Context) ([]*entity.Album, error)
	Update(ctx context.Context, album *entity.Album) error
	Delete(ctx context.Context, id int64) error
	AddMedia(ctx context.Context, albumID, mediaID int64) error
	RemoveMedia(ctx context.Context, albumID, mediaID int64) error
}

type albumRepository struct {
	db *bun.DB
}

func NewAlbumRepository(db *bun.DB) AlbumRepository {
	return &albumRepository{db: db}
}

func (r *albumRepository) Create(ctx context.Context, album *entity.Album) (*entity.Album, error) {
	_, err := r.db.NewInsert().Model(album).Returning("*").Exec(ctx)
	if err != nil {
		return nil, err
	}
	return album, nil
}

func (r *albumRepository) GetByID(ctx context.Context, id int64) (*entity.Album, error) {
	album := new(entity.Album)

	err := r.db.NewSelect().Model(album).Where("a.id = ?", id).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperr.ErrNotFound
		}
		return nil, err
	}

	return album, nil
}

func (r *albumRepository) List(ctx context.Context) ([]*entity.Album, error) {
	var albums []*entity.Album

	err := r.db.NewSelect().Model(&albums).Order("a.id ASC").Scan(ctx)
	if err != nil {
		return nil, err
	}

	return albums, nil
}

func (r *albumRepository) Update(ctx context.Context, album *entity.Album) error {
	res, err := r.db.NewUpdate().Model(album).Column("name").WherePK().Exec(ctx)
	if err != nil {
		return err
	}
	return ensureAffected(res)
}

func (r *albumRepository) Delete(ctx context.Context, id int64) error {
	// No FK cascade on SQLite here, so remove the link rows and the album in one
	// transaction to avoid orphaned album_media entries.
	return r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewDelete().
			Model((*entity.AlbumMedia)(nil)).
			Where("album_id = ?", id).
			Exec(ctx); err != nil {
			return err
		}

		res, err := tx.NewDelete().
			Model((*entity.Album)(nil)).
			Where("id = ?", id).
			Exec(ctx)
		if err != nil {
			return err
		}
		return ensureAffected(res)
	})
}

func (r *albumRepository) AddMedia(ctx context.Context, albumID, mediaID int64) error {
	link := &entity.AlbumMedia{AlbumID: albumID, MediaID: mediaID}

	_, err := r.db.NewInsert().Model(link).Exec(ctx)
	return err
}

func (r *albumRepository) RemoveMedia(ctx context.Context, albumID, mediaID int64) error {
	res, err := r.db.NewDelete().
		Model((*entity.AlbumMedia)(nil)).
		Where("album_id = ? AND media_id = ?", albumID, mediaID).
		Exec(ctx)
	if err != nil {
		return err
	}
	return ensureAffected(res)
}

// ensureAffected maps a zero-row write to a not-found error.
func ensureAffected(res sql.Result) error {
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return apperr.ErrNotFound
	}
	return nil
}
