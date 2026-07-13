package controller

import (
	"fhs/app/internal/dto"
	"fhs/app/internal/service"

	"github.com/go-fuego/fuego"
)

const groupName = "/media"

type MediaController struct {
	srv service.MediaService
}

func NewMediaController(srv service.MediaService) *MediaController {
	return &MediaController{srv: srv}
}

func (c *MediaController) RegisterRoutes(server *fuego.Server) {
	router := fuego.Group(server, groupName)

	fuego.Get(router, "/{id}", c.getMedia)
	fuego.Post(router, "/", c.createMedia)
}

func (c *MediaController) getMedia(fc fuego.ContextNoBody) (*dto.MediaDto, error) {
	id, err := fc.PathParamIntErr("id")
	if err != nil {
		return nil, fuego.BadRequestError{Title: "Invalid ID", Err: err}
	}

	media, err := c.srv.GetMedia(fc.Context(), int64(id))
	if err != nil {
		return nil, fuego.NotFoundError{Title: "Media not found", Err: err}
	}

	return media, nil
}

func (c *MediaController) createMedia(fc fuego.ContextWithBody[dto.CreateMediaDto]) (*dto.MediaDto, error) {
	dto, err := fc.Body()
	if err != nil {
		return nil, err
	}

	return c.srv.CreateMedia(fc.Context(), &dto)
}
