package handlers

import (
	"AABBCCDD/app/db"
	"AABBCCDD/app/views/landing"

	"github.com/anthdm/gothkit/pkg/kit"
)

func HandleLandingIndex(kit *kit.Kit) error {
	db.Query.NewSelect().Scan(kit.Request.Context())
	return kit.Render(landing.Index())
}
