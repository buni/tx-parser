package dto

type GetCurrentBlockRequest struct{}

type GetCurrentBlockResponse struct {
	Height string `json:"height"`
}
