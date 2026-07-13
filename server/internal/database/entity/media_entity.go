package entity

import "github.com/uptrace/bun"

type Media struct {
	bun.BaseModel `bun:"table:media,alias:m"`

	ID    int64  `bun:"id,pk,autoincrement"`
	Title string `bun:"title"`
	URL   string `bun:"url"`
	Path  string `bun:"path"`

	Albums []*Album `bun:"m2m:album_media,join:Media=Album"`
}
