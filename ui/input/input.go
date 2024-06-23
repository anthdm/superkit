package input

import (
	"github.com/a-h/templ"

	"github.com/anthdm/superkit/ui"
)

const defaultInputClass = "flex h-10 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"

func Input(opts ...func(*templ.Attributes)) templ.Attributes {
	return ui.CreateAttrs(defaultInputClass, "", opts...)
}
