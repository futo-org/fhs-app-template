package dto

type CreateAlbumDto struct {
	Name string `json:"name" validate:"required"`
}

type UpdateAlbumDto struct {
	Name string `json:"name" validate:"required"`
}

type AddMediaToAlbumDto struct {
	MediaID int64 `json:"mediaId" validate:"required"`
}

type AlbumDto struct {
	ID    int64      `json:"id"`
	Name  string     `json:"name"`
	Media []MediaDto `json:"media"`
}
