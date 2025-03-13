package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TxManager interface {
	BeginTx(ctx context.Context) 				(pgx.Tx, error)
	CommitTx(ctx context.Context, tx pgx.Tx) 	error
	RollbackTx(ctx context.Context, tx pgx.Tx) 	error
}

type transactionManager struct {
	pool *pgxpool.Pool
}

func NewTransactionManager(pool *pgxpool.Pool) TxManager {
	return &transactionManager{
		pool: pool,
	}
}

func (tm *transactionManager) BeginTx(ctx context.Context) (pgx.Tx, error) {
	tx, err := tm.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (tm *transactionManager) CommitTx(ctx context.Context, tx pgx.Tx) error {
	return tx.Commit(ctx)
}

func (tm *transactionManager) RollbackTx(ctx context.Context, tx pgx.Tx) error {
	return tx.Rollback(ctx)
}