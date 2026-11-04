package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/arvaliullin/shortener/internal/api/http/dto"
	"github.com/arvaliullin/shortener/internal/core/domain"
	"github.com/go-chi/chi/v5"
)

// Update меняет исходный адрес короткой ссылки.
func (h *URLHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	var req dto.ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || !req.IsValid() {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	originalURL := strings.TrimSpace(req.URL)
	err := h.urlService.Update(r.Context(), id, originalURL)
	if errors.Is(err, domain.ErrURLNotFound) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto.NewURLItem(h.baseURL, id, originalURL))
}
