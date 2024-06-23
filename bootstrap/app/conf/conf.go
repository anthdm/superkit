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
	SUPERKIT_ENV:     v.Env[string]("SUPERKIT_ENV", v.Rules(v.In([]string{"development", "staging", "production"})), "development"),
	HTTP_LISTEN_ADDR: v.Env[string]("HTTP_LISTEN_ADDR", v.Rules(v.Min(3)), ":3000"),
	SUPERKIT_SECRET:  v.Env[string]("SUPERKIT_SECRET", v.Required, v.Rules(v.Min(32))),
}
