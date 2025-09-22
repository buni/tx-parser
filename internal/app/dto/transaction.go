package dto

import "github.com/buni/tx-parser/internal/app/entity"

type ListAddressTransactionsRequest struct {
	Address   string           `json:"-" in:"path=address" validate:"required"`
	TokenType entity.TokenType `json:"-" in:"path=tokenType" validate:"required"`
}

type ListAddressTransactionsResponse struct {
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	TokenType   string `json:"token_type"`
	ToAddress   string `json:"to_address"`
	FromAddress string `json:"from_address"`
	Address     string `json:"address"`
	BlockNumber string `json:"block_number"`
	Hash        string `json:"hash"`
	Value       string `json:"value"`
}

func TransactionsFromEntities(entityTxs []entity.Transaction) []Transaction {
	txs := make([]Transaction, 0, len(entityTxs))
	for k := range entityTxs {
		txs = append(txs, TransactionFromEntity(entityTxs[k]))
	}
	return txs
}

func TransactionFromEntity(tx entity.Transaction) Transaction {
	return Transaction{
		ID:          tx.ID,
		TokenType:   tx.TokenType.String(),
		ToAddress:   tx.ToAddress,
		FromAddress: tx.FromAddress,
		Address:     tx.Address,
		Hash:        tx.Hash,
		Value:       tx.Value, // this is the value in wei for ETH
		UserID:      tx.UserID,
		BlockNumber: tx.BlockNumber,
	}
}
