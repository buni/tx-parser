package entity

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"github.com/buni/tx-parser/internal/pkg/pubsub"
)

const (
	TokenTypeInvliad TokenType = iota
	TokenTypeETH
)

const (
	TransactionTypeInvalid TransactionType = iota
	TransactionTypeDebit
	TransactionTypeCredit
)

const (
	TransactionCreatedTopic = "transaction.created"
)

type TransactionType int

//go:generate go tool enumer -type=TokenType,TransactionType -trimprefix=TokenType,TransactionType -transform=snake -output=transaction_enum.go -json -sql -text

type TokenType int

type Transaction struct {
	ID              string          `db:"id"`
	TokenType       TokenType       `db:"token_type"`
	BlockNumber     string          `db:"block_hash"`
	TransactionType TransactionType `db:"transaction_type"`
	UserID          string          `db:"user_id"`
	Address         string          `db:"address"`
	ToAddress       string          `db:"to_address"`
	FromAddress     string          `db:"from_address"`
	Hash            string          `db:"hash"`
	Value           string          `db:"value"`
	CreatedAt       time.Time       `db:"created_at"`
}

func NewTransaction(tokenType TokenType, transactionType TransactionType, userID, address, toAddress, fromAddress, hash, value, blockNumber string) (Transaction, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return Transaction{}, fmt.Errorf("generate uuid: %w", err)
	}

	return Transaction{
		ID:              id.String(),
		TokenType:       tokenType,
		TransactionType: transactionType,
		UserID:          userID,
		Address:         address,
		ToAddress:       toAddress,
		FromAddress:     fromAddress,
		Hash:            hash,
		Value:           value,
		BlockNumber:     blockNumber,
		CreatedAt:       time.Now().UTC().Truncate(time.Millisecond),
	}, nil
}

func NewTransactionCreatedEvent(tx Transaction) (*pubsub.Message, error) {
	msg, err := pubsub.NewJSONMessage(tx, nil)
	if err != nil {
		return nil, fmt.Errorf("json message: %w", err)
	}

	msg.Topic = TransactionCreatedTopic

	return msg, nil
}
