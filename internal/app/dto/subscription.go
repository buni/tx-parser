package dto

type SubscribeRequest struct {
	Address string `json:"address" validate:"required"`
}

type SubscribeResponse struct{}
