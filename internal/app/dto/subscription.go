package dto

import "github.com/buni/tx-parser/internal/app/entity"

type SubscribeRequest struct {
	Address   string           `json:"address" validate:"required"`
	UserID    string           `json:"user_id" validate:"required"`
	TokenType entity.TokenType `json:"-" in:"path=tokenType" validate:"required"`
}

type SubscribeResponse struct{}
