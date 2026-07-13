package controller

import (
	"errors"
	"net/http"

	"fhs/app/internal/apperr"
	"fhs/app/internal/dto"
	"fhs/app/internal/service"

	"github.com/go-fuego/fuego"
)

const albumGroupName = "/albums"

type AlbumController struct {
	srv service.AlbumService
}

func NewAlbumController(srv service.AlbumService) *AlbumController {
	return &AlbumController{srv: srv}
}

func (c *AlbumController) RegisterRoutes(server *fuego.Server) {
	router := fuego.Group(server, albumGroupName)

	fuego.Post(router, "/", c.createAlbum, fuego.OptionDefaultStatusCode(http.StatusCreated))
	fuego.Get(router, "/", c.listAlbums)
	fuego.Get(router, "/{id}", c.getAlbum)
	fuego.Put(router, "/{id}", c.updateAlbum)
	fuego.Delete(router, "/{id}", c.deleteAlbum, fuego.OptionDefaultStatusCode(http.StatusNoContent))
	fuego.Post(router, "/{id}/media", c.addMedia, fuego.OptionDefaultStatusCode(http.StatusNoContent))
	fuego.Delete(router, "/{id}/media/{mediaId}", c.removeMedia, fuego.OptionDefaultStatusCode(http.StatusNoContent))
}

func (c *AlbumController) createAlbum(fc fuego.ContextWithBody[dto.CreateAlbumDto]) (*dto.AlbumDto, error) {
	body, err := fc.Body()
	if err != nil {
		return nil, err
	}

	album, err := c.srv.CreateAlbum(fc.Context(), &body)
	if err != nil {
		return nil, toHTTPError(err)
	}
	return album, nil // 201 Created
}

func (c *AlbumController) listAlbums(fc fuego.ContextNoBody) ([]*dto.AlbumDto, error) {
	albums, err := c.srv.ListAlbums(fc.Context())
	if err != nil {
		return nil, toHTTPError(err)
	}
	return albums, nil
}

func (c *AlbumController) getAlbum(fc fuego.ContextNoBody) (*dto.AlbumDto, error) {
	id, err := fc.PathParamIntErr("id")
	if err != nil {
		return nil, fuego.BadRequestError{Title: "Invalid ID", Err: err}
	}

	album, err := c.srv.GetAlbum(fc.Context(), int64(id))
	if err != nil {
		return nil, toHTTPError(err)
	}
	return album, nil
}

func (c *AlbumController) updateAlbum(fc fuego.ContextWithBody[dto.UpdateAlbumDto]) (*dto.AlbumDto, error) {
	id, err := fc.PathParamIntErr("id")
	if err != nil {
		return nil, fuego.BadRequestError{Title: "Invalid ID", Err: err}
	}

	body, err := fc.Body()
	if err != nil {
		return nil, err
	}

	album, err := c.srv.UpdateAlbum(fc.Context(), int64(id), &body)
	if err != nil {
		return nil, toHTTPError(err)
	}
	return album, nil
}

func (c *AlbumController) deleteAlbum(fc fuego.ContextNoBody) (any, error) {
	id, err := fc.PathParamIntErr("id")
	if err != nil {
		return nil, fuego.BadRequestError{Title: "Invalid ID", Err: err}
	}

	if err := c.srv.DeleteAlbum(fc.Context(), int64(id)); err != nil {
		return nil, toHTTPError(err)
	}
	return nil, nil
}

func (c *AlbumController) addMedia(fc fuego.ContextWithBody[dto.AddMediaToAlbumDto]) (any, error) {
	id, err := fc.PathParamIntErr("id")
	if err != nil {
		return nil, fuego.BadRequestError{Title: "Invalid ID", Err: err}
	}

	body, err := fc.Body()
	if err != nil {
		return nil, err
	}

	if err := c.srv.AddMedia(fc.Context(), int64(id), body.MediaID); err != nil {
		return nil, toHTTPError(err)
	}
	return nil, nil
}

func (c *AlbumController) removeMedia(fc fuego.ContextNoBody) (any, error) {
	id, err := fc.PathParamIntErr("id")
	if err != nil {
		return nil, fuego.BadRequestError{Title: "Invalid ID", Err: err}
	}

	mediaID, err := fc.PathParamIntErr("mediaId")
	if err != nil {
		return nil, fuego.BadRequestError{Title: "Invalid media ID", Err: err}
	}

	if err := c.srv.RemoveMedia(fc.Context(), int64(id), int64(mediaID)); err != nil {
		return nil, toHTTPError(err)
	}
	return nil, nil
}

func toHTTPError(err error) error {
	switch {
	case errors.Is(err, apperr.ErrNotFound):
		return fuego.NotFoundError{Title: "Not found", Err: err}
	case errors.Is(err, apperr.ErrInvalidInput):
		return fuego.BadRequestError{Title: "Invalid input", Err: err}
	default:
		return err
	}
}
