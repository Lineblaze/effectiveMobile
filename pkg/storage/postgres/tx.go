package postgres

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Tx struct {
	db pgx.Tx
}

type TxReq func(tx Tx) error

type TxRunner interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

func ExecTx(ctx context.Context, runner TxRunner, req TxReq) error {
	pgxTx, err := runner.Begin(ctx)
	if err != nil {
		return err
	}

	tx := Tx{
		db: pgxTx,
	}

	defer tx.Rollback(context.TODO())

	err = req(tx)
	if err != nil {
		return err
	}
	return tx.Commit(context.TODO())
}

func (p Tx) Stats() *pgxpool.Stat {
	return nil
}

func (p Tx) Begin(ctx context.Context) (pgx.Tx, error) {
	return p.db.Begin(ctx)
}

func (p Tx) Rollback(ctx context.Context) {
	_ = p.db.Rollback(ctx)
}

func (p Tx) Commit(ctx context.Context) error {
	return p.db.Commit(ctx)
}

func (p Tx) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	return p.db.Query(ctx, query, args...)
}

func (p Tx) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	rows, err := p.db.Query(ctx, query, args...)
	if err != nil {
		return err
	}
	return pgxscan.ScanOne(dest, rows)
}

func (p Tx) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	rows, err := p.db.Query(ctx, query, args...)
	if err != nil {
		return err
	}
	return pgxscan.ScanAll(dest, rows)
}

func (p Tx) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	return p.Exec(ctx, sql, arguments...)
}

func (p Tx) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return p.db.QueryRow(ctx, query, args...)
}
