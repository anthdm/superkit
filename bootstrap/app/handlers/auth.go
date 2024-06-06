package handlers

import (
	"net/http"

	"example-app/app/types"

	"github.com/anthdm/gothkit/pkg/kit"
)

func HandleAuthentication(w http.ResponseWriter, r *http.Request) (kit.Auth, error) {
	return types.AuthUser{}, nil
}
