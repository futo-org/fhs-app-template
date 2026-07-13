package entity

import "github.com/uptrace/bun"

type Album struct {
	bun.BaseModel `bun:"table:albums,alias:a"`

	ID   int64  `bun:"id,pk,autoincrement"`
	Name string `bun:"name,notnull"`

	Media []Media `bun:"m2m:album_media,join:Album=Media"`
}
