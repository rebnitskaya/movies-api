package middleware

import (
	"log"
	"net/http"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Println("panic recovered:", rec)
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
