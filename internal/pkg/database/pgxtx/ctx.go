package pgxtx

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type txKey struct{}

func SetTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func GetTx(ctx context.Context) pgx.Tx {
	tx, _ := ctx.Value(txKey{}).(pgx.Tx) //nolint:revive
	return tx
}

func GetTxOrCreate(ctx context.Context, db *pgxpool.Pool, opts pgx.TxOptions) (context.Context, pgx.Tx, bool, error) {
	tx := GetTx(ctx)
	if tx != nil {
		return ctx, tx, false, nil
	}

	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return nil, tx, false, fmt.Errorf("failed to create transaction: %w", err)
	}

	return SetTx(ctx, tx), tx, true, nil
}
