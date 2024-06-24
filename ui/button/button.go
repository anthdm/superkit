package button

import (
	"github.com/anthdm/superkit/ui"

	"github.com/a-h/templ"
)

const (
	buttonBaseClass          = "inline-flex items-center justify-center px-4 py-2 font-medium text-sm tracking-wide transition-colors duration-200 rounded-md focus:ring focus:shadow-outline focus:outline-none"
	buttonVariantPrimary     = "text-primary-foreground bg-primary focus:ring-primary hover:bg-primary/90"
	buttonVariantOutline     = "text-primary border border-primary hover:bg-secondary focus:ring-primary"
	buttonVariantSecondary   = "text-primary bg-secondary hover:bg-secondary/80"
	buttonVariantDestructive = "text-primary bg-destructive hover:bg-destructive/80"
)

func New(opts ...func(*templ.Attributes)) templ.Attributes {
	return ui.CreateAttrs(buttonBaseClass, buttonVariantPrimary, opts...)
}

func Outline(opts ...func(*templ.Attributes)) templ.Attributes {
	return appendVariant("outline", opts...)
}

func Primary(opts ...func(*templ.Attributes)) templ.Attributes {
	return appendVariant("primary", opts...)
}

func Secondary(opts ...func(*templ.Attributes)) templ.Attributes {
	return appendVariant("secondary", opts...)
}

func Destructive(opts ...func(*templ.Attributes)) templ.Attributes {
	return appendVariant("destructive", opts...)
}

func Variant(variant string) func(*templ.Attributes) {
	return func(attrs *templ.Attributes) {
		att := *attrs
		switch variant {
		case "primary":
			att["class"] = ui.Merge(buttonBaseClass, buttonVariantPrimary)
		case "outline":
			att["class"] = ui.Merge(buttonBaseClass, buttonVariantOutline)
		case "secondary":
			att["class"] = ui.Merge(buttonBaseClass, buttonVariantSecondary)
		case "destructive":
			att["class"] = ui.Merge(buttonBaseClass, buttonVariantDestructive)
		}
	}
}

func appendVariant(variant string, opts ...func(*templ.Attributes)) templ.Attributes {
	opt := []func(*templ.Attributes){
		Variant(variant),
	}
	opt = append(opt, opts...)
	return New(opt...)
}
