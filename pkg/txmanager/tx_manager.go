package txmanager

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5"
)

type TxManager interface {
	Run(ctx context.Context, fn func(ctx context.Context) error) error
	Tx(ctx context.Context) pgx.Tx
}


type CtxKey struct {}


type pgxTxManager struct {
	pool *pgxpool.Pool
}

func NewPgxTxManager(pool *pgxpool.Pool) TxManager {
	return &pgxTxManager{
		pool: pool,
	}
}

func (p *pgxTxManager) Run(ctx context.Context, fn func(ctx context.Context) error) error {

	tx, err := p.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	ctxTx := context.WithValue(context.Background(), CtxKey{}, tx)

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback(ctx)
			panic(r)
		} else if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()


	if err := fn(ctxTx); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (p *pgxTxManager) Tx(ctx context.Context) pgx.Tx {
	tx, ok := ctx.Value(CtxKey{}).(pgx.Tx)
	if !ok {
		tx, _ := p.pool.BeginTx(ctx, pgx.TxOptions{})
		return tx
	}
	return tx
}