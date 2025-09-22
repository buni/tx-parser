package entity

import (
	"math/big"
	"time"
)

type Block struct {
	TokenType TokenType `db:"token_type"`
	Height    big.Int   `db:"height"`
	Hash      string    `db:"hash"`
	CreatedAt time.Time `db:"created_at"`
}
