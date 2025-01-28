package handler

import (
	"fmt"
	"net/http"

	"github.com/buni/tx-parser/internal/app/contract"
	"github.com/buni/tx-parser/internal/app/dto"
	"github.com/buni/tx-parser/internal/pkg/handler"
	"github.com/go-chi/chi/v5"
)

type SubscriptionHandler struct {
	svc contract.SubscriptionService
}

func NewSubscriptionHandler(svc contract.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{
		svc: svc,
	}
}

func (h *SubscriptionHandler) Subscribe(w http.ResponseWriter, req *http.Request, reqBody *dto.SubscribeRequest) (*dto.SubscribeResponse, error) {
	err := h.svc.Subscribe(req.Context(), reqBody)
	if err != nil {
		return nil, fmt.Errorf("subscribe: %w", err)
	}

	w.WriteHeader(http.StatusCreated)

	return &dto.SubscribeResponse{}, nil
}

func (h *SubscriptionHandler) RegisterRoutes(r chi.Router) {
	r.Post("/addresses/subscribe", handler.WrapDefault(h.Subscribe))
}
