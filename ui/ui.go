package ui

import (
	"fmt"

	"github.com/a-h/templ"
)

func CreateAttrs(baseClass string, defaultClass string, opts ...func(*templ.Attributes)) templ.Attributes {
	attrs := templ.Attributes{
		"class": baseClass + " " + defaultClass,
	}
	for _, o := range opts {
		o(&attrs)
	}
	return attrs
}

func Merge(a, b string) string {
	return fmt.Sprintf("%s %s", a, b)
}

func Class(class string) func(*templ.Attributes) {
	return func(attrs *templ.Attributes) {
		attr := *attrs
		class := attr["class"].(string) + " " + class
		attr["class"] = class
	}
}
