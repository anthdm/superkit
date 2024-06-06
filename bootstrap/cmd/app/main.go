package main

import (
	"example-app/app"
	"fmt"
	"net/http"
	"os"

	"github.com/anthdm/gothkit/pkg/kit"
	"github.com/go-chi/chi/v5"
)

func main() {
	router := chi.NewMux()

	app.InitializeMiddleware(router)

	if kit.IsDevelopment() {
		router.Handle("/public/*", disableCache(staticDev()))
	}

	kit.UseErrorHandler(app.ErrorHandler)
	router.HandleFunc("/*", kit.Handler(app.NotFoundHandler))

	app.InitializeRoutes(router)

	fmt.Printf("application running in %s at %s\n", kit.Env(), "http://localhost:7331")

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
