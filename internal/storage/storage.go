package storage

import (
	"nevermore/internal/storage/postgres"
)

//go:generate mockery --name=Storage --dir=. --output=./mocks
type Storage interface {
	DB() postgres.Repo
}

type repo struct {
	psql postgres.Repo
}

func (r *repo) DB() postgres.Repo {
	return r.psql
}

func New(pcfg postgres.Config) (Storage, error) {
	psql, err := postgres.NewDB(pcfg)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	result := &repo{
		psql: psql,
	}
	return result, nil
}
