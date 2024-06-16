package auth

import (
	"github.com/anthdm/superkit/kit"
	"github.com/go-chi/chi/v5"
)

func InitializeRoutes(router chi.Router) {
	authConfig := kit.AuthenticationConfig{
		AuthFunc:    AuthenticateUser,
		RedirectURL: "/login",
	}

	router.Group(func(auth chi.Router) {
		auth.Use(kit.WithAuthentication(authConfig, false))
		auth.Get("/login", kit.Handler(HandleAuthIndex))
		auth.Post("/login", kit.Handler(HandleAuthCreate))
		auth.Delete("/logout", kit.Handler(HandleAuthDelete))

		auth.Get("/signup", kit.Handler(HandleSignupIndex))
		auth.Post("/signup", kit.Handler(HandleSignupCreate))
	})

	router.Group(func(auth chi.Router) {
		auth.Use(kit.WithAuthentication(authConfig, true))
		auth.Get("/profile", kit.Handler(HandleProfileShow))
		auth.Put("/profile", kit.Handler(HandleProfileUpdate))
	})
}
