package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"example-app/db"

	"github.com/go-chi/chi/v5"
)

func main() {
	db, err := db.New()
	if err != nil {
		log.Fatal(err)
	}

	_ = db

	// Routes configuration
	router := chi.NewMux()
	if true {
		router.Handle("/*", disableCache(staticDev()))
	}
	initializeRoutes(router, db)

	listenAddr := os.Getenv("HTTP_LISTEN_ADDR")
	slog.Info("application started", "listenAddr", listenAddr)

	http.ListenAndServe(os.Getenv("HTTP_LISTEN_ADDR"), router)
}

func staticDev() http.Handler {
	return http.StripPrefix("/public/", http.FileServerFS(os.DirFS("public")))
}

func disableCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}
