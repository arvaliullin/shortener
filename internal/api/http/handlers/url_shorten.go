package handlers

import (
	"io"
	"net/http"
	"strings"
)

func (h *URLHandler) Shorten(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	originalURL := strings.TrimSpace(string(body))
	if originalURL == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	id, err := h.urlService.Shorten(r.Context(), originalURL)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(strings.TrimRight(h.baseURL, "/") + "/" + id))
}
