package server

import "github.com/gin-gonic/gin"

func NewServer() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	e := gin.New()
	e.Use(gin.Logger())

	return e
}
