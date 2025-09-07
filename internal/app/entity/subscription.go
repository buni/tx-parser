package entity

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

type Subscription struct {
	ID        string    `db:"id"`
	TokenType TokenType `db:"token_type"`
	UserID    string    `db:"user_id"`
	Address   string    `db:"address"`
	CreatedAt time.Time `db:"created_at"`
}

func NewSubscription(tokenType TokenType, userID, address string) (Subscription, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return Subscription{}, fmt.Errorf("generate uuid: %w", err)
	}

	return Subscription{
		ID:        id.String(),
		TokenType: tokenType,
		UserID:    userID,
		Address:   address,
		CreatedAt: time.Now().UTC().Truncate(time.Millisecond),
	}, nil
}
