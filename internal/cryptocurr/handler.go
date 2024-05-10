package cryptocurr

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type response struct {
	Data []DTO `json:"data"`
}

type service interface {
	GetAll(ctx context.Context) ([]DTO, error)
	GetByTitle(ctx context.Context, title string) (*DTO, error)
}

type Handler struct {
	logger  *slog.Logger
	service service
}

func NewHandler(logger *slog.Logger, service service) *Handler {
	return &Handler{
		logger:  logger,
		service: service,
	}
}

func (h *Handler) GetAll(ctx *gin.Context) {
	dtos, err := h.service.GetAll(ctx)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)

		return
	}

	resp := response{Data: make([]DTO, 0, cap(dtos))}
	resp.Data = append(resp.Data, dtos...)
	ctx.IndentedJSON(http.StatusOK, resp)
}

func (h *Handler) GetByTitle(ctx *gin.Context) {
	title := ctx.Param("title")
	dto, err := h.service.GetByTitle(ctx, title)

	switch {
	case err == nil:
		resp := response{Data: []DTO{*dto}}
		ctx.IndentedJSON(http.StatusOK, resp)

	case errors.Is(err, errNotFound):
		ctx.Status(http.StatusNotFound)

	case err != nil:
		ctx.Status(http.StatusInternalServerError)
	}
}
