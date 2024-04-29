package rest

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/pavelk123/cryptocurrency-service/internal/entity"
	"log/slog"
)

type service interface {
	GetAll(ctx context.Context) ([]*entity.DTO, error)
	GetByTitle(ctx context.Context, title string) (*entity.DTO, error)
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
		ctx.JSON(500, err.Error())
	}

	ctx.JSON(200, []interface{}{dtos})
}

func (h *Handler) GetByTitle(ctx *gin.Context) {
	title := ctx.Param("title")

	dto, err := h.service.GetByTitle(ctx, title)
	if err != nil {
		h.logger.Error(err.Error())
	}

	ctx.JSON(200, dto)
}
