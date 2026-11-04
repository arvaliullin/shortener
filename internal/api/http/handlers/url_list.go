package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/arvaliullin/shortener/internal/api/http/dto"
)

// List возвращает сохраненные короткие ссылки.
func (h *URLHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.urlService.List(r.Context())
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	response := make([]dto.URLItem, 0, len(items))
	for _, item := range items {
		response = append(response, dto.NewURLItem(h.baseURL, item.ID, item.Original))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
