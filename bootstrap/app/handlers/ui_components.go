package handlers

import (
	uicomponents "AABBCCDD/app/views/ui_components"

	"github.com/anthdm/superkit/kit"
)

func HandleUIComponentsIndex(kit *kit.Kit) error {
	return kit.Render(uicomponents.Index())
}
