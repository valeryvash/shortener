package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RedirectPathParams struct {
	ShortAddr string `uri:"shortAddr" binding:"required,len=8"`
}

func (h *Handler) RedirectHandlerGIN(c *gin.Context) {
	var (
		redirectPathParams RedirectPathParams
		fullAddress        string
		err                error
	)

	if err = c.ShouldBindUri(&redirectPathParams); err != nil {
		errBadRequest(c)
		return
	}

	if fullAddress, err = h.service.GetFullAddr(redirectPathParams.ShortAddr); err != nil {
		errBadRequest(c)
		return
	}

	// c.Header("Location", fullAddress)
	// c.Status(http.StatusTemporaryRedirect)
	//or just
	c.Redirect(http.StatusTemporaryRedirect, fullAddress)
}
