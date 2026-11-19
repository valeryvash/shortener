package repository_test

import (
	"github.com/valeryvash/shortener/internal/repository"
	"github.com/valeryvash/shortener/internal/service"
)

var _ service.Repository = (*repository.Repository)(nil)
