package handlers

import (
	"net/http"

	"AABBCCDD/app/types"

	"github.com/anthdm/gothkit/pkg/kit"
)

func HandleAuthentication(w http.ResponseWriter, r *http.Request) (kit.Auth, error) {
	return types.AuthUser{}, nil
}
