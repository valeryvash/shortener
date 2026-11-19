package handler

import (
	"net/http"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

func (h *Handler) RedirectHandlerGIN(c *gin.Context) {
	shortAddr := c.Param("shortAddr")
	if utf8.RuneCountInString(shortAddr) != 8 {
		errBadRequest(c)
		return
	}

	fullAddr, err := h.service.GetFullAddr(shortAddr)
	if err != nil {
		errBadRequest(c)
		return
	}

	// c.Header("Location", fullAddr)
	// c.String(http.StatusTemporaryRedirect, "")
	// or just
	c.Redirect(http.StatusTemporaryRedirect, fullAddr)
}
