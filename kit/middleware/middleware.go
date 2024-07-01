package middleware

import (
	"context"
	"net/http"
)

type (
	RequestKey         struct{}
	ResponseHeadersKey struct{}
)

func WithRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), RequestKey{}, r)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
