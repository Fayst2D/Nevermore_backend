package service

import (
	"nevermore/internal/service/user"
	"nevermore/internal/storage"

	"github.com/gammazero/workerpool"

	"nevermore/pkg/hash"
)

type Service interface {
	User() user.Service
}

type service struct {
	user user.Service
}

func New(st storage.Storage,
	hash hash.PasswordHasher,
	wp *workerpool.WorkerPool) Service {

	result := &service{
		user: user.New(st),
	}

	return result
}

func (s *service) User() user.Service {
	return s.user
}
