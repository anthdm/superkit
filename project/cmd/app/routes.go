package main

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/anthdm/gothkit/pkg/kit"
	"github.com/go-chi/chi/v5"

	"example-app/handlers"
	"example-app/views/errors"
)

// Define your routes in here
func initializeRoutes(router *chi.Mux, db *sql.DB) {
	// Configure the error handler
	kit.UseErrorHandler(func(kit *kit.Kit, err error) {
		slog.Error("internal server error", "err", err.Error(), "path", kit.Request.URL.Path)
		kit.Render(errors.Error500())
	})

	// Comment out to configure your authentication
	authConfig := kit.AuthenticationConfig{
		AuthFunc:    handleAuthentication,
		RedirectURL: "/login",
	}

	landingHandler := handlers.NewLandingHandler(db)
	// Routes that "might" have an authenticated user
	router.Group(func(app chi.Router) {
		app.Use(kit.WithAuthentication(authConfig, false)) // strict set to false

		// Routes
		app.Get("/", kit.Handler(landingHandler.HandleIndex))
	})

	// Routes that "must" have an authenticated user or else they
	// will be redirected to the configured redirectURL, set in the
	// AuthenticationConfig.
	router.Group(func(app chi.Router) {
		app.Use(kit.WithAuthentication(authConfig, true)) // strict set to true

		// Routes
		// app.Get("/path", kit.Handler(myHandler.HandleIndex))
	})
}

type AuthUser struct {
	ID       int
	Email    string
	LoggedIn bool
}

func (user AuthUser) Check() bool {
	return user.ID > 0 && user.LoggedIn
}

func handleAuthentication(w http.ResponseWriter, r *http.Request) (kit.Auth, error) {
	return AuthUser{}, nil
}
