package pgxtx

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Query interface {
	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)
}

type QueryRow interface {
	QueryRow(ctx context.Context, query string, args ...any) pgx.Row
}

type Exec interface {
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
}

type Querier interface {
	Query
	QueryRow
	Exec
}

type TxWrapper struct {
	*pgxpool.Pool
	txOptions pgx.TxOptions
}

// NewTxWrapper returns a new TxWrapper.
func NewTxWrapper(db *pgxpool.Pool, txOptions pgx.TxOptions) *TxWrapper {
	return &TxWrapper{
		Pool:      db,
		txOptions: txOptions,
	}
}

// Exec...
func (db *TxWrapper) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	var tx Exec
	tx = db.Pool

	ctxTx := GetTx(ctx)
	if ctxTx != nil {
		tx = ctxTx
	}

	return tx.Exec(ctx, query, args...) //nolint:wrapcheck
}

// Query ...
func (db *TxWrapper) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	var tx Query
	tx = db.Pool

	ctxTx := GetTx(ctx)
	if ctxTx != nil {
		tx = ctxTx
	}

	return tx.Query(ctx, query, args...) //nolint:wrapcheck
}

// QueryRow ...
func (db *TxWrapper) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	var tx QueryRow
	tx = db.Pool

	ctxTx := GetTx(ctx)
	if ctxTx != nil {
		tx = ctxTx
	}

	return tx.QueryRow(ctx, query, args...) //nolint:wrapcheck
}
