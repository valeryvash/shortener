package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateHandlerGIN(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		errBadRequest(c)
		return
	}

	shortAddr, err := h.service.CreateShortAddr(string(body))
	if err != nil {
		errBadRequest(c)
		return
	}
	
	c.Header("Content-Type", "text/plain")
	c.String(http.StatusCreated, "http://localhost:8080/"+string(shortAddr))
}