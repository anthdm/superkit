package kit

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

var store *sessions.CookieStore

type HandlerFunc func(kit *Kit) error

type ErrorHandlerFunc func(kit *Kit, err error)

type AuthKey struct{}

type Auth interface {
	Check() bool
}

var (
	errorHandler = func(kit *Kit, err error) {
		kit.Text(http.StatusInternalServerError, err.Error())
	}
)

type DefaultAuth struct{}

func (DefaultAuth) Check() bool { return false }

type Kit struct {
	Response http.ResponseWriter
	Request  *http.Request
}

func UseErrorHandler(h ErrorHandlerFunc) { errorHandler = h }

func (kit *Kit) Auth() Auth {
	value, ok := kit.Request.Context().Value(AuthKey{}).(Auth)
	if !ok {
		slog.Warn("kit authentication not set")
		return DefaultAuth{}
	}
	return value
}

// GetSession return a session by its name. GetSession always
// returns a session even if it does not exist.
func (kit *Kit) GetSession(name string) *sessions.Session {
	sess, _ := store.Get(kit.Request, name)
	return sess
}

// Redirect with HTMX support.
func (kit *Kit) Redirect(status int, url string) error {
	if len(kit.Request.Header.Get("HX-Request")) > 0 {
		kit.Response.Header().Set("HX-Redirect", url)
		kit.Response.WriteHeader(http.StatusSeeOther)
		return nil
	}
	http.Redirect(kit.Response, kit.Request, url, status)
	return nil
}

func (kit *Kit) FormValue(name string) string {
	return kit.Request.PostFormValue(name)
}

func (kit *Kit) JSON(status int, v any) error {
	kit.Response.WriteHeader(status)
	kit.Response.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(kit.Response).Encode(v)
}

func (kit *Kit) Text(status int, msg string) error {
	kit.Response.WriteHeader(status)
	kit.Response.Header().Set("Content-Type", "text/plain")
	_, err := kit.Response.Write([]byte(msg))
	return err
}

func (kit *Kit) Bytes(status int, b []byte) error {
	kit.Response.WriteHeader(status)
	kit.Response.Header().Set("Content-Type", "text/plain")
	_, err := kit.Response.Write(b)
	return err
}

func (kit *Kit) Render(c templ.Component) error {
	return c.Render(kit.Request.Context(), kit.Response)
}

func (kit *Kit) Getenv(name string, def string) string {
	return Getenv(name, def)
}

func Handler(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		kit := &Kit{
			Response: w,
			Request:  r,
		}
		if err := h(kit); err != nil {
			if errorHandler != nil {
				errorHandler(kit, err)
				return
			}
			kit.Text(http.StatusInternalServerError, err.Error())
		}
	}
}

type AuthenticationConfig struct {
	AuthFunc    func(*Kit) (Auth, error)
	RedirectURL string
}

func WithAuthentication(config AuthenticationConfig, strict bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			kit := &Kit{
				Response: w,
				Request:  r,
			}
			auth, err := config.AuthFunc(kit)
			if err != nil {
				errorHandler(kit, err)
				return
			}
			if strict && !auth.Check() && r.URL.Path != config.RedirectURL {
				kit.Redirect(http.StatusSeeOther, config.RedirectURL)
				return
			}
			ctx := context.WithValue(r.Context(), AuthKey{}, auth)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Getenv(name string, def string) string {
	env := os.Getenv(name)
	if len(env) == 0 {
		return def
	}
	return env
}

func IsDevelopment() bool {
	return os.Getenv("SUPERKIT_ENV") == "development"
}

func IsProduction() bool {
	return os.Getenv("SUPERKIT_ENV") == "production"
}

func Env() string {
	return os.Getenv("SUPERKIT_ENV")
}

// initialize the store here so the environment variables are
// already initialized. Calling NewCookieStore() from outside of
// a function scope won't work.
func Setup() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	appSecret := os.Getenv("SUPERKIT_SECRET")
	if len(appSecret) < 32 {
		// For security reasons we are calling os.Exit(1) here so Go's panic recover won't
		// recover the application without a valid SUPERKIT_SECRET set.
		fmt.Println("invalid SUPERKIT_SECRET variable. Are you sure you have set the SUPERKIT_SECRET in your .env file?")
		os.Exit(1)
	}
	store = sessions.NewCookieStore([]byte(appSecret))
}
