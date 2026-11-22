package repository

import "sync"

type Repository struct {
	mu sync.Mutex
	storage map[string]string
}

func NewRepository() *Repository {
	return &Repository{
		storage: make(map[string]string),
	}
}
