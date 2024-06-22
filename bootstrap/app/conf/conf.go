package conf

import (
	v "github.com/anthdm/superkit/validate"
)

// Application config

var Env = struct {
	SUPERKIT_ENV     string
	HTTP_LISTEN_ADDR string
	SUPERKIT_SECRET  string
}{
	SUPERKIT_ENV:     v.Env("SUPERKIT_ENV", v.In([]string{"development", "staging", "production"})).Default("development").Validate(),
	HTTP_LISTEN_ADDR: v.Env("HTTP_LISTEN_ADDR", v.Min(3)).Default(":3000").Validate(),
	SUPERKIT_SECRET:  v.Env("SUPERKIT_SECRET", v.Required, v.Min(32)).Validate(),
}
