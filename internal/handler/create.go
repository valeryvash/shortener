package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateHandlerGIN(c *gin.Context) {
	var (
		body      string
		shortAddr string
		err       error
	)

	if err = c.ShouldBindPlain(&body); err != nil {
		errBadRequest(c)
		return
	}

	if shortAddr, err = h.service.CreateShortAddr(body); err != nil {
		errBadRequest(c)
		return
	}

	c.String(
		http.StatusCreated,
		"%s/%s", h.redirectAddr, shortAddr,
	)
}
