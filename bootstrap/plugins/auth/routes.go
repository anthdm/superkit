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

	router.Get("/email/verify", kit.Handler(HandleEmailVerify))
	router.Post("/resend-email-verification", kit.Handler(HandleResendVerificationCode))

	router.Group(func(auth chi.Router) {
		auth.Use(kit.WithAuthentication(authConfig, false))
		auth.Get("/login", kit.Handler(HandleLoginIndex))
		auth.Post("/login", kit.Handler(HandleLoginCreate))
		auth.Delete("/logout", kit.Handler(HandleLoginDelete))

		auth.Get("/signup", kit.Handler(HandleSignupIndex))
		auth.Post("/signup", kit.Handler(HandleSignupCreate))

	})

	router.Group(func(auth chi.Router) {
		auth.Use(kit.WithAuthentication(authConfig, true))
		auth.Get("/profile", kit.Handler(HandleProfileShow))
		auth.Put("/profile", kit.Handler(HandleProfileUpdate))
	})
}
