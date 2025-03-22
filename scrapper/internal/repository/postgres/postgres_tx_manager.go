package postgres

import (
	"context"
	"errors"
	"log/slog"

	txmanager "github.com/Ranik23/tbank-tech/scrapper/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type txPostgresManager struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

type txKey struct{}

func NewTxManager(pool *pgxpool.Pool, logger *slog.Logger) txmanager.TxManager {
	return &txPostgresManager{
		pool:   pool,
		logger: logger,
	}
}

func (t *txPostgresManager) GetExecutor(ctx context.Context) txmanager.Executor {
	tx, ok := ctx.Value(txKey{}).(txmanager.Executor)
	if ok {
		return tx
	}
	return t.pool
}

func (t *txPostgresManager) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	opts := pgx.TxOptions{
		IsoLevel:   pgx.Serializable,
		AccessMode: pgx.ReadWrite,
	}

	var err error

	tx, err := t.pool.BeginTx(ctx, opts)
	if err != nil {
		t.logger.Error("Failed to start the transaction", slog.String("error", err.Error()))
		return err
	}

	defer func() {
		if err != nil {
			if rollbackError := tx.Rollback(ctx); rollbackError != nil && !errors.Is(rollbackError, pgx.ErrTxClosed) {
				t.logger.Error("Failed to rollback tx", slog.String("error", rollbackError.Error()))
			}
			t.logger.Info("Rollback successfully done!")
		}
	}()

	ctx = context.WithValue(ctx, txKey{}, tx)

	err = fn(ctx)
	if err != nil {
		t.logger.Error("Failed to run the chain", slog.String("error", err.Error()))
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		t.logger.Error("Failed to commit the transaction", slog.String("error", err.Error()))
		return err
	}

	return nil
}
