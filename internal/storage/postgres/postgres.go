package postgres

import (
	"context"
	"fmt"
	"nevermore/internal/storage/postgres/author"
	"nevermore/internal/storage/postgres/book"
	"nevermore/internal/storage/postgres/saved_author"
	"nevermore/internal/storage/postgres/user"

	"github.com/jmoiron/sqlx"
)

type Repo interface {
	BeginTx(ctx context.Context) (*sqlx.Tx, error)
	User() user.Repo
	Author() author.Repo
	SavedAuthor() saved_author.Repo
	Book() book.Repo
}

type repo struct {
	db          *sqlx.DB
	user        user.Repo
	author      author.Repo
	savedAuthor saved_author.Repo
	book        book.Repo
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
		db:          db,
		user:        user.New(db),
		author:      author.New(db),
		savedAuthor: saved_author.New(db),
		book:        book.New(db),
	}
	return result, nil
}

func (r *repo) User() user.Repo {
	return r.user
}
func (r *repo) Author() author.Repo {
	return r.author
}
func (r *repo) SavedAuthor() saved_author.Repo {
	return r.savedAuthor
}
func (r *repo) Book() book.Repo {
	return r.book
}
