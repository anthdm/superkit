package middleware

import (
	"context"
	"net/http"
)

type RequestURLKey struct{}

func WithRequestURL(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), RequestURLKey{}, r.URL)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
