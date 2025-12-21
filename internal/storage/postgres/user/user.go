package user

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"nevermore/internal/dto"
	model "nevermore/internal/model/user"
)

type Repo interface {
	Create(ctx context.Context, user *model.User) error
	Get(ctx context.Context, id int) (*dto.UserGetResponse, error)
	GetById(ctx context.Context, id int) (model.User, error)
	Update(ctx context.Context, u model.User) error
	Delete(ctx context.Context, id int) error
	GetByEmail(ctx context.Context, email string) (model.User, error)
	GetByVerificationToken(ctx context.Context, token string) (model.User, error)
	VerifyEmail(ctx context.Context, token string) error
	UpdateVerificationToken(ctx context.Context, email string, token string) error
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

func (r *repo) Create(ctx context.Context, user *model.User) error {
	query := `insert into users 
				(name, phone_number, email, password, role, photo, created_at, email_verified, verification_token) 
			  values ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.ExecContext(
		ctx,
		query,
		user.Name,
		user.PhoneNumber,
		user.Email,
		user.Password,
		user.Role,
		user.Photo,
		user.CreatedAt,
		user.IsEmailVerified,
		user.VerificationToken,
	)

	return err
}

func (r *repo) Update(ctx context.Context, user model.User) error {
	query := `update users
          set name = $1, 
              phone_number = $2, 
              password = $3,
              photo = $4,
          where id = $5`

	_, err := r.db.ExecContext(
		ctx,
		query,
		user.Name,
		user.PhoneNumber,
		user.Password,
		user.Photo,
		user.Id,
	)

	return err
}

func (r *repo) Delete(ctx context.Context, id int) error {
	query := "update users set deleted_at = $1 where id = $2 and deleted_at is null"

	deletedAt := time.Now().UTC()

	_, err := r.db.ExecContext(ctx, query, deletedAt, id)

	return err
}

func (r *repo) Get(ctx context.Context, id int) (*dto.UserGetResponse, error) {
	var user dto.UserGetResponse

	query := "select name, phone_number, photo, email, role from users where id = $1 and deleted_at is null"

	err := r.db.GetContext(ctx, &user, query, id)
	return &user, err
}

func (r *repo) GetById(ctx context.Context, id int) (model.User, error) {
	var user model.User

	query := "select id, name, phone_number, photo, email, password, role, email_verified, verification_token, created_at, deleted_at from users where id = $1 and deleted_at is null"

	err := r.db.GetContext(ctx, &user, query, id)
	return user, err
}

func (r *repo) GetByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User

	query := "select id, name, phone_number, photo, email, password, role, email_verified, verification_token, created_at, deleted_at from users where email = $1 and deleted_at is null"

	err := r.db.GetContext(ctx, &user, query, email)

	return user, err
}

func (r *repo) GetByVerificationToken(ctx context.Context, token string) (model.User, error) {
	var user model.User

	query := "select id, name, phone_number, photo, email, password, role, email_verified, verification_token, created_at, deleted_at from users where verification_token = $1 and deleted_at is null"

	err := r.db.GetContext(ctx, &user, query, token)

	return user, err
}

func (r *repo) VerifyEmail(ctx context.Context, token string) error {
	query := "update users set email_verified = true, verification_token = null where verification_token = $1 and deleted_at is null"

	result, err := r.db.ExecContext(ctx, query, token)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("verification token not found or already used")
	}

	return nil
}

func (r *repo) UpdateVerificationToken(ctx context.Context, email string, token string) error {
	query := "update users set verification_token = $1, email_verified = false where email = $2 and deleted_at is null"

	_, err := r.db.ExecContext(ctx, query, token, email)
	return err
}
