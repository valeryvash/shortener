package main

import (
	"log"
	"net/http"

	"github.com/valeryvash/shortener/internal/handler"
	"github.com/valeryvash/shortener/internal/repository"
	"github.com/valeryvash/shortener/internal/service"
)


func main() {
	repo := repository.NewRepository()
	service := service.NewService(repo)
	handler := handler.NewHandler(service)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
