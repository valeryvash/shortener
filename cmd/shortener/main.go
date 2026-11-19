package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/valeryvash/shortener/internal/handler"
	"github.com/valeryvash/shortener/internal/repository"
	"github.com/valeryvash/shortener/internal/service"
)

func main() {

	repo := repository.NewRepository()
	service := service.NewService(repo)
	handler := handler.NewHandler(service)

	e := gin.Default()
	handler.RegisterRoutesGIN(e)

	if err := e.Run("localhost:8080"); err != nil {
		log.Fatal(err)
	}
}
