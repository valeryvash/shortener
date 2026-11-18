package handler

import (
	"fmt"
	"io"
	"net/http"
)

func (h *Handler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shortAddr, err := h.service.CreateShortAddr(string(body))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, _ = fmt.Fprintf(w, "http://localhost:8080/%s", string(shortAddr))
}