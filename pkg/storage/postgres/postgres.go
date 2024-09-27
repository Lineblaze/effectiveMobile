package postgres

import (
	"context"
	"effectiveMobile/config"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

const (
	_defaultConnectionAttempts = 10
	_defaultConnectionTimeout  = time.Second
	_maxConnections            = int32(800)
	_minConnections            = int32(50)
	_maxConnectionLifeTime     = time.Second * 300
	_maxIdleLifeTime           = time.Second * 15
)

type Postgres interface {
	Stats() *pgxpool.Stat
	Query(query string, args ...any) (pgx.Rows, error)
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	Exec(query string, args ...any) (pgconn.CommandTag, error)
	QueryRow(query string, args ...interface{}) pgx.Row
	TxRunner
}

type Pool struct {
	db *pgxpool.Pool
}

func InitPsqlDB(c *config.Config) (Postgres, error) {
	connectionUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Postgres.Host,
		c.Postgres.Port,
		c.Postgres.User,
		c.Postgres.Password,
		c.Postgres.DBName,
		c.Postgres.SSLMode)
	connectionUrl += fmt.Sprintf(" pool_max_conns=%d pool_min_conns=%d pool_max_conn_lifetime=%v pool_max_conn_idle_time=%v",
		_maxConnections, _minConnections, _maxConnectionLifeTime, _maxIdleLifeTime)
	connectionAttempts := _defaultConnectionAttempts
	var result *pgxpool.Pool
	var err error
	for connectionAttempts > 0 {
		result, err = pgxpool.New(context.Background(), connectionUrl)
		if err == nil {
			break
		}

		log.Printf("ATTEMPT %d TO CONNECT TO POSTGRES BY URL %s FAILED: %s\n", connectionAttempts, connectionUrl, err.Error())

		connectionAttempts--

		time.Sleep(_defaultConnectionTimeout)
	}

	if result == nil {
		log.Printf("POSTGRES CONNECTION(%s) ERROR: %s\n", connectionUrl, err.Error())
		return nil, err
	}

	return &Pool{db: result}, nil
}

func (p Pool) Stats() *pgxpool.Stat {
	return p.db.Stat()
}

func (p Pool) Begin(ctx context.Context) (pgx.Tx, error) {
	return p.db.Begin(ctx)
}

func (p Pool) Query(query string, args ...any) (pgx.Rows, error) {
	return p.db.Query(context.Background(), query, args...)
}

func (p Pool) Get(dest interface{}, query string, args ...interface{}) error {
	rows, err := p.db.Query(context.Background(), query, args...)
	if err != nil {
		return err
	}
	return pgxscan.ScanOne(dest, rows)
}

func (p Pool) Select(dest interface{}, query string, args ...interface{}) error {
	rows, err := p.db.Query(context.Background(), query, args...)
	if err != nil {
		return err
	}
	return pgxscan.ScanAll(dest, rows)
}

func (p Pool) Exec(query string, args ...interface{}) (pgconn.CommandTag, error) {
	return p.db.Exec(context.Background(), query, args...)
}

func (p Pool) QueryRow(query string, args ...interface{}) pgx.Row {
	return p.db.QueryRow(context.Background(), query, args...)
}
