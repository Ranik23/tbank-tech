package txmanager

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TxManager interface {
	BeginTx(ctx context.Context) 				(pgx.Tx, error)
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
	return tm.pool.Begin(ctx)
}
