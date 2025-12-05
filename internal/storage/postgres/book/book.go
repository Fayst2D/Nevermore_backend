package book

import (
	"context"
	"nevermore/internal/dto"

	"github.com/jmoiron/sqlx"
)

type Repo interface {
	Create(ctx context.Context, tx *sqlx.Tx, req *dto.CreateBookRequest) (int, error)
	SaveFirstPage(ctx context.Context, tx *sqlx.Tx, url string, bookId int) error
	GetByAuthor(ctx context.Context, authorID int) ([]dto.GetBookRequest, error)
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

func (r *repo) Create(ctx context.Context, tx *sqlx.Tx, req *dto.CreateBookRequest) (int, error) {
	query := `insert into books 
				(title, description, author, uploaded_by, status, url) 
			  values ($1, $2, $3, $4, $5, $6) returning id`

	var id int

	err := tx.QueryRowxContext(ctx,
		query,
		req.Title,
		req.Description,
		req.Author,
		req.UploadedBy,
		0,
		req.FileUrl,
	).Scan(&id)

	return id, err
}

func (r *repo) SaveFirstPage(ctx context.Context, tx *sqlx.Tx, url string, bookId int) error {
	query := `update books set url = $1 where id = $2`

	_, err := tx.ExecContext(ctx, query, url, bookId)

	return err
}

func (r *repo) GetByAuthor(ctx context.Context, authorID int) ([]dto.GetBookRequest, error) {
	var result []*dto.GetBookRequest

	query := `SELECT id, title, description, author, uploaded_by, status, url, created_at 
              FROM books 
              WHERE author = $1 
              ORDER BY created_at DESC`

	var books []dto.GetBookRequest
	err := r.db.SelectContext(ctx, &result, query, authorID)

	return books, err
}
