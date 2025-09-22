package dto

import "github.com/buni/tx-parser/internal/app/entity"

type GetCurrentBlockRequest struct {
	TokenType entity.TokenType `json:"-" in:"path=tokenType" validate:"required"`
}

type GetCurrentBlockResponse struct {
	Height string `json:"height"`
}
