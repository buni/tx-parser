package dto

import "github.com/buni/tx-parser/internal/app/entity"

type ListAddressTransactionsRequest struct {
	Address string `json:"-" in:"path=address" validate:"required"`
}

type ListAddressTransactionsResponse struct {
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	ID        string `json:"id"`
	TokenType string `json:"token_type"`
	To        string `json:"to"`
	From      string `json:"from"`
	Address   string `json:"address"`
	Hash      string `json:"hash"`
	Value     string `json:"value"`
}

func TransactionsFromEntities(entityTxs []entity.Transaction) []Transaction {
	txs := make([]Transaction, 0, len(entityTxs))
	for _, v := range entityTxs {
		txs = append(txs, TransactionFromEntity(v))
	}
	return txs
}

func TransactionFromEntity(tx entity.Transaction) Transaction {
	return Transaction{
		ID:        tx.ID,
		TokenType: tx.TokenType.String(),
		To:        tx.To,
		From:      tx.From,
		Address:   tx.Address,
		Hash:      tx.Hash,
		Value:     tx.Value,
	}
}
