package postgres

import (
	"context"
	"fmt"
	"nevermore/internal/storage/postgres/user"

	"github.com/jmoiron/sqlx"
)

type Repo interface {
	BeginTx(ctx context.Context) (*sqlx.Tx, error)
	User() user.Repo
}

type repo struct {
	db   *sqlx.DB
	user user.Repo
}

func (r *repo) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	return r.db.BeginTxx(ctx, nil)
}

func NewDB(cfg Config) (Repo, error) {
	db, err := sqlx.Connect(cfg.Driver, cfg.URL)
	fmt.Println("CONNECTED TO PSQL")
	if err != nil {
		return nil, err
	}

	result := &repo{
		db:   db,
		user: user.New(db),
	}
	return result, nil
}

func (r *repo) User() user.Repo {
	return r.user
}
