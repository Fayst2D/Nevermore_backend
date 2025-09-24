package user

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"nevermore/internal/dto"
	model "nevermore/internal/model/user"
)

type Repo interface {
	Create(ctx context.Context, user *model.User) error
	Get(ctx context.Context, id int) (*dto.UserGetResponse, error)
	Update(ctx context.Context, u model.User) error
	Delete(ctx context.Context, id int) error
	GetByEmail(ctx context.Context, email string) (model.User, error)
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
				(name, phone_number, email, password, role, photo, created_at) 
			  values ($1, $2, $3, $4, $5, $6, $7)`

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
	)

	return err
}

func (r *repo) Update(ctx context.Context, user model.User) error {
	query := `update users
          set name = $1, 
              phone_number = $2, 
              email = $3, 
              password = $4,
              photo = $5,
              role = $6
          where id = $7`

	_, err := r.db.ExecContext(
		ctx,
		query,
		user.Name,
		user.PhoneNumber,
		user.Email,
		user.Password,
		user.Photo,
		user.Role,
		user.Id,
	)

	return err
}

func (r *repo) Delete(ctx context.Context, id int) error {
	query := "update users set deleted_at = $1 where id = $2 and deleted_at is null"

	deletedAt := time.Now()

	_, err := r.db.ExecContext(ctx, query, deletedAt, id)

	return err
}

func (r *repo) Get(ctx context.Context, id int) (*dto.UserGetResponse, error) {
	var user dto.UserGetResponse

	query := "select name, phone_number, photo, email, role from users where id = $1 and deleted_at is null"

	err := r.db.GetContext(ctx, &user, query, id)
	return &user, err
}

func (r *repo) GetByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User

	query := "select id, name, phone_number, photo, email, password, role, created_at, deleted_at from users where email = $1 and deleted_at is null"

	err := r.db.GetContext(ctx, &user, query, email)

	return user, err
}
