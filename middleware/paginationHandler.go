package middleware

import (
	"fmt"
	"net/http"
	"strconv"
)

func GetPaginationParams(r *http.Request) (int, int, error) {
	page := 1
	limit := 5

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		value, err := strconv.Atoi(pageStr)
		if err != nil || value < 1 {
			return 0, 0, fmt.Errorf("Invalid page")
		}
		page = value
	}

	if limitStr := r.URL.Query().Get("size"); limitStr != "" {
		value, err := strconv.Atoi(limitStr)
		if err != nil || value < 1 {
			return 0, 0, fmt.Errorf("Invalid limit")
		}
		limit = value
	}

	return page, limit, nil
}
