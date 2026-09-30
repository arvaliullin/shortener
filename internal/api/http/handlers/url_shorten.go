package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/arvaliullin/shortener/internal/api/http/dto"
)

// Shorten обрабатывает запрос на создание короткой ссылки
func (h *URLHandler) Shorten(w http.ResponseWriter, r *http.Request) {
	var req dto.ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if !req.IsValid() {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	id, err := h.urlService.Shorten(r.Context(), strings.TrimSpace(req.URL))
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.NewShortenResponse(h.baseURL, id))
}
