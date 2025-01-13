package utils

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if status == http.StatusNoContent {
		return
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("Failed to write JSON response.", slog.String("error", err.Error()))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
