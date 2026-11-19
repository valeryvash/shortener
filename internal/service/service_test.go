package service_test

import (
	"github.com/valeryvash/shortener/internal/handler"
	"github.com/valeryvash/shortener/internal/service"
)

var _ handler.Service = (*service.Service)(nil)
