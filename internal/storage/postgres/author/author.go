package author

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"nevermore/internal/dto"
	model "nevermore/internal/model/author"
)

type Repo interface {
	Create(ctx context.Context, author *model.Author) error
	Get(ctx context.Context, id int) (*dto.AuthorGetResponse, error)
	Update(ctx context.Context, u model.Author) error
	Delete(ctx context.Context, id int) error
}

type repo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) Repo {
	result := &repo{
		db: db,
	}

	return result
}

func (r *repo) Create(ctx context.Context, author *model.Author) error {
	query := `insert into authors 
				(name, biography, photo_url, created_at) 
			  values ($1, $2, $3, $4)`

	_, err := r.db.ExecContext(
		ctx,
		query,
		author.Name,
		author.Biography,
		author.Photo,
		author.CreatedAt,
	)

	return err
}

func (r *repo) Update(ctx context.Context, author model.Author) error {
	query := `update authors
          set name = $1, 
              biography = $2, 
              photo_url = $3, 
          where id = $4`

	_, err := r.db.ExecContext(
		ctx,
		query,
		author.Name,
		author.Biography,
		author.Photo,
		author.Id,
	)

	return err
}

func (r *repo) Delete(ctx context.Context, id int) error {
	query := "update authors set deleted_at = $1 where id = $2"

	deletedAt := time.Now().UTC()

	_, err := r.db.ExecContext(ctx, query, deletedAt, id)

	return err
}

func (r *repo) Get(ctx context.Context, id int) (*dto.AuthorGetResponse, error) {
	var author dto.AuthorGetResponse

	query := "select name, biography, photo_url from authors where id = $1"

	err := r.db.GetContext(ctx, &author, query, id)
	return &author, err
}
