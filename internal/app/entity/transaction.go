package entity

const (
	TokenTypeInvliad TokenType = iota
	TokenTypeETH
)

//go:generate enumer -type=TokenType -trimprefix=TokenType -transform=snake -output=transaction_enum.go -json -sql -text

type TokenType int

type Transaction struct {
	ID        string    `json:"id"`
	TokenType TokenType `json:"token_type"`
	To        string    `json:"to"`
	From      string    `json:"from"`
	Address   string    `json:"address"`
	Hash      string    `json:"hash"`
	Value     string    `json:"value"`
}
