module AABBCCDD

go 1.22.4

// uncomment for local development on the superkit core.
// replace github.com/anthdm/superkit => ../superkit

require (
	github.com/a-h/templ v0.2.731
	github.com/anthdm/superkit v0.0.0-20240622052611-30be5bb82e0d
	github.com/go-chi/chi/v5 v5.0.14
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/mattn/go-sqlite3 v1.14.22
	golang.org/x/crypto v0.24.0
	gorm.io/driver/sqlite v1.5.6
	gorm.io/gorm v1.25.10
)

require (
	github.com/gorilla/securecookie v1.1.2 // indirect
	github.com/gorilla/sessions v1.3.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
)
