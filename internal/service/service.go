package service

import (
	"nevermore/internal/service/author"
	"nevermore/internal/service/saved_author"
	"nevermore/internal/service/user"
	"nevermore/internal/storage"

	"github.com/gammazero/workerpool"

	"nevermore/pkg/hash"
)

type Service interface {
	User() user.Service
	Author() author.Service
	SavedAuthor() saved_author.Service
}

type service struct {
	user        user.Service
	author      author.Service
	savedAuthor saved_author.Service
}

func New(st storage.Storage,
	hash hash.PasswordHasher,
	wp *workerpool.WorkerPool) Service {

	result := &service{
		user:   user.New(st),
		author: author.New(st),
	}

	return result
}

func (s *service) User() user.Service {
	return s.user
}
func (s *service) Author() author.Service {
	return s.author
}
func (s *service) SavedAuthor() saved_author.Service {
	return s.savedAuthor
}
