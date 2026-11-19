package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:generate mockery --name=Service --output=mocks --outpkg=mocks
type Service interface {
	CreateShortAddr(fullAddr string) (shortAddr string, err error)
	GetFullAddr(shortAddr string) (fullAddr string, err error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutesGIN(e *gin.Engine) {
	e.POST("/", h.CreateHandlerGIN)
	e.GET("/:shortAddr", h.RedirectHandlerGIN)

	e.NoRoute(errBadRequest)
	e.NoMethod(errBadRequest)
}

func errBadRequest(c *gin.Context) {
	c.String(http.StatusBadRequest, "")
}
