package model

import (
	"context"
	"database/sql"
)

func Transaction(
	ctx context.Context,
	db *sql.DB,
	txOptions *sql.TxOptions,
	f func(ctx context.Context, tx *sql.Tx) error,
) error {
	tx, err := db.BeginTx(ctx, txOptions)
	if err != nil {
		return err
	}
	if err := f(ctx, tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
