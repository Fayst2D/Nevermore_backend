package storage

import (
	minio2 "nevermore/internal/storage/minio"
	"nevermore/internal/storage/postgres"
	"nevermore/internal/storage/redis"
)

//go:generate mockery --name=Storage --dir=. --output=./mocks
type Storage interface {
	DB() postgres.Repo
	Cash() redis.Storage
	Cloud() minio2.Storage
}

type repo struct {
	rds    redis.Storage
	psql   postgres.Repo
	photos minio2.Storage
}

func (r *repo) DB() postgres.Repo {
	return r.psql
}

func (r *repo) Cash() redis.Storage {
	return r.rds
}

func (r *repo) Cloud() minio2.Storage {
	return r.photos
}

func New(pcfg postgres.Config, rdscfg redis.Config, photoes minio2.Config) (Storage, error) {
	psql, err := postgres.NewDB(pcfg)
	if err != nil {
		return nil, err
	}

	rds, err := redis.New(rdscfg)
	if err != nil {
		return nil, err
	}

	cloud, err := minio2.New(photoes)
	if err != nil {
		return nil, err
	}

	result := &repo{
		psql:   psql,
		rds:    rds,
		photos: cloud,
	}
	return result, nil
}
