package handler

import (
	"context"
	"fmt"

	"github.com/buni/tx-parser/internal/app/contract"
	"github.com/buni/tx-parser/internal/app/dto"
	"github.com/buni/tx-parser/internal/pkg/handler"
	"github.com/go-chi/chi/v5"
)

type TransactionHandler struct {
	svc contract.TransactionService
}

func NewTransactionHandler(svc contract.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		svc: svc,
	}
}

func (h *TransactionHandler) ListAddressTransactions(ctx context.Context, req *dto.ListAddressTransactionsRequest) (*dto.ListAddressTransactionsResponse, error) {
	txs, err := h.svc.ListAddressTransactions(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("list address transactions: %w", err)
	}

	return &dto.ListAddressTransactionsResponse{
		Transactions: dto.TransactionsFromEntities(txs),
	}, nil
}

func (h *TransactionHandler) RegisterRoutes(r chi.Router) {
	r.Get("/addresses/{address}/transactions", handler.WrapDefaultBasic(h.ListAddressTransactions))
}
