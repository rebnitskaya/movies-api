package handlers

import (
	"encoding/json"
	"net/http"
)

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		return
	}

	err := json.NewEncoder(w).Encode("Hi!")
	if err != nil {
		http.Error(w, "Some server error during root processing. Try again later!", http.StatusInternalServerError)
	}
}
