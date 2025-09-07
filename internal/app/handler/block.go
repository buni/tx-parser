package handler

import (
	"context"
	"fmt"

	"github.com/go-chi/chi/v5"

	"github.com/buni/tx-parser/internal/app/contract"
	"github.com/buni/tx-parser/internal/app/dto"
	"github.com/buni/tx-parser/internal/pkg/handler"
)

// BlockHandler is the handler for block-related operations.
type BlockHandler struct {
	svc contract.BlockService
}

// NewBlockHandler creates a new BlockHandler.
func NewBlockHandler(svc contract.BlockService) *BlockHandler {
	return &BlockHandler{
		svc: svc,
	}
}

func (h *BlockHandler) GetCurrentBlock(ctx context.Context, req *dto.GetCurrentBlockRequest) (*dto.GetCurrentBlockResponse, error) {
	height, err := h.svc.GetCurrentBlock(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("get current block: %w", err)
	}

	return &dto.GetCurrentBlockResponse{
		Height: height,
	}, nil
}

func (h *BlockHandler) RegisterRoutes(r chi.Router) {
	r.Get("/blocks/current", handler.WrapDefaultBasic(h.GetCurrentBlock))
}
