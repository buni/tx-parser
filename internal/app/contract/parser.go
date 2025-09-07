package contract

import (
	"context"
)

//go:generate go tool mockgen -source=parser.go -destination=mock/parser_mocks.go -package contract_mock

type ParserService interface {
	ParseNextBlock(ctx context.Context) error
	ParseBlock(ctx context.Context, height string) (err error)
}
