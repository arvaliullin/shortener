package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Redirect обрабатывает запрос на перенаправление по ID
func (h *URLHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	originalURL, err := h.urlService.Resolve(r.Context(), id)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
}
