package handlers

import (
	"AABBCCDD/app/views/showcase"

	"github.com/anthdm/superkit/kit"
)

func HandleShowcaseIndex(kit *kit.Kit) error {
	return kit.Render(showcase.Index())
}
