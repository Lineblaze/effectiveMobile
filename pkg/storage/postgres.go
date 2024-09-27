package storage

import (
	"effectiveMobile/config"
	"fmt"
	"github.com/guregu/null/v5"
	"github.com/jmoiron/sqlx"
	"time"
)

type Comment struct {
	Id              int64    `db:"id"`
	UserId          int      `db:"user_id"`
	CommentId       null.Int `db:"comment_id"`
	Content         string   `db:"content"`
	Level           int
	VoteCount       int     `db:"voteCount"`
	UpdatedAt       float64 `db:"updated_at"`
	UpdatedAtNormal string
}

func InitPsqlDB(c *config.Config) (*sqlx.DB, error) {
	connectionUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Postgres.Host,
		c.Postgres.Port,
		c.Postgres.User,
		c.Postgres.Password,
		c.Postgres.DBName,
		c.Postgres.SSLMode)
	database, err := sqlx.Connect(c.Postgres.PgDriver, connectionUrl)
	if err != nil {
		return nil, err
	}
	database.DB.SetConnMaxIdleTime(60 * time.Second)
	database.DB.SetMaxOpenConns(300)
	database.DB.SetMaxIdleConns(300)
	database.DB.SetConnMaxLifetime(60 * time.Second)
	database.DB.SetConnMaxIdleTime(60 * time.Second)
	return database, nil
}
