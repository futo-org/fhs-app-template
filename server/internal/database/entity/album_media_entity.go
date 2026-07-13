package entity

import "github.com/uptrace/bun"

type AlbumMedia struct {
	bun.BaseModel `bun:"table:album_media,alias:am"`

	AlbumID int64 `bun:"album_id,notnull"`
	MediaID int64 `bun:"media_id,notnull"`

	Album *Album `bun:"rel:belongs-to,join:album_id=id"`
	Media *Media `bun:"rel:belongs-to,join:media_id=id"`
}
