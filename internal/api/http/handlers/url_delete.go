package handlers

import (
	"errors"
	"net/http"

	"github.com/arvaliullin/shortener/internal/core/domain"
	"github.com/go-chi/chi/v5"
)

// Delete удаляет короткую ссылку.
func (h *URLHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	err := h.urlService.Delete(r.Context(), id)
	if errors.Is(err, domain.ErrURLNotFound) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
