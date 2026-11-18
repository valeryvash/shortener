package handler

import (
	"net/http"
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

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", h.ServeHTTP)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/" && r.Method == http.MethodPost:
		h.CreateHandler(w, r)
	case r.URL.Path != "/" && r.Method == http.MethodGet:
		h.RedirectHandler(w, r)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}
