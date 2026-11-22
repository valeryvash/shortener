package main

import (
	"log"
	"os"

	"github.com/valeryvash/shortener/internal/config"
	"github.com/valeryvash/shortener/internal/handler"
	"github.com/valeryvash/shortener/internal/repository"
	"github.com/valeryvash/shortener/internal/server"
	"github.com/valeryvash/shortener/internal/service"
)

func main() {
	config, err := config.NewConfig(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewRepository()
	service := service.NewService(repo)
	handler := handler.NewHandler(service, config.RedirectAddress)

	e := server.NewServer()
	handler.RegisterRoutesGIN(e)

	if err := e.Run(config.BindAddress); err != nil {
		log.Fatal(err)
	}
}
