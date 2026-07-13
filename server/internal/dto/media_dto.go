package dto

type CreateMediaDto struct {
	Name string `json:"name" validate:"required"`
	Type string `json:"type" validate:"required"`
	Path string `json:"path" validate:"required"`
}

type MediaDto struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
	Path  string `json:"path"`
}
