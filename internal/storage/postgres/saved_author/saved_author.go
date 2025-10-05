package saved_author

import (
	"context"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"nevermore/internal/dto"
	model "nevermore/internal/model/saved_author"
)

type Repo interface {
	Create(ctx context.Context, author *model.SavedAuthor) error
	Delete(ctx context.Context, id int) error
	GetSavedAuthorsList(ctx context.Context, userID int) ([]*dto.AuthorGetResponse, error)
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

func (r *repo) Create(ctx context.Context, author *model.SavedAuthor) error {
	query := `insert into saved_authors 
				(user_id, author_id, created_at) 
			  values ($1, $2, $3, $4)`

	_, err := r.db.ExecContext(
		ctx,
		query,
		author.UserId,
		author.AuthorId,
		author.CreatedAt,
	)

	return err
}

func (r *repo) Delete(ctx context.Context, id int) error {
	query := "delete from saved_authors where id = $1"

	_, err := r.db.ExecContext(ctx, query, id)

	return err
}

func (r *repo) GetSavedAuthorsList(ctx context.Context, userID int) ([]*dto.AuthorGetResponse, error) {
	var authors []*dto.AuthorGetResponse

	query := `SELECT authors.name, authors.biography, authors.photo_url, saved_authors.created_at 
              FROM saved_authors 
              JOIN authors ON saved_authors.author_id = authors.id 
              WHERE saved_authors.user_id = ? 
              ORDER BY saved_authors.created_at DESC`

	err := r.db.SelectContext(ctx, &authors, query, userID)
	if err != nil {
		return nil, err
	}

	return authors, nil
}
