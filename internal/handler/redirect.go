package handler

import "net/http"

func (h *Handler) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	path := []rune(r.URL.Path)
	if len(path) != 9 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	path = path[1:]
	fullAddr, err := h.service.GetFullAddr(string(path))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Add("Location", fullAddr)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
