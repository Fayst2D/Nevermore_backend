package storage

import (
	"nevermore/internal/storage/minio"
	"nevermore/internal/storage/postgres"
)

//go:generate mockery --name=Storage --dir=. --output=./mocks
type Storage interface {
	DB() postgres.Repo
	Cloud() minio.Storage
}

type repo struct {
	psql    postgres.Repo
	photoes minio.Storage
}

func (r *repo) DB() postgres.Repo {
	return r.psql
}

func New(pcfg postgres.Config, photoes minio.Config) (Storage, error) {
	psql, err := postgres.NewDB(pcfg)
	if err != nil {
		return nil, err
	}

	cloud, err := minio.New(photoes)
	if err != nil {
		return nil, err
	}

	result := &repo{
		psql:    psql,
		photoes: cloud,
	}
	return result, nil
}

func (r *repo) Cloud() minio.Storage {
	return r.photoes
}
