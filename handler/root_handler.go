package handler

import (
	"encoding/json"
	"net/http"
)

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		return
	}

	err := json.NewEncoder(w).Encode("Try another request!")
	if err != nil {
		http.Error(w, "Some server error during root processing. Try again later!", http.StatusInternalServerError)
	}
}
