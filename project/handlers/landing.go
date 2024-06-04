package handlers

import (
	"database/sql"
	"example-app/views/landing"

	"github.com/anthdm/gothkit/pkg/kit"
)

type LandingHandler struct {
	db *sql.DB
}

func NewLandingHandler(db *sql.DB) *LandingHandler {
	return &LandingHandler{
		db: db,
	}
}

func (h *LandingHandler) HandleIndex(kit *kit.Kit) error {
	return kit.Render(landing.Index())
}
