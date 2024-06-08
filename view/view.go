package view

import (
	"context"
	"fmt"
	"net/url"

	"github.com/a-h/templ"
	"github.com/anthdm/gothkit/kit"
	"github.com/anthdm/gothkit/kit/middleware"
)

func Asset(name string) templ.SafeURL {
	return templ.URL(fmt.Sprintf("/public/assets/%s", name))
}

// Auth is a view helper function that returns the current Auth.
// If Auth is not set a default auth will be returned
func Auth(ctx context.Context) kit.Auth {
	value, ok := ctx.Value(kit.AuthKey{}).(kit.Auth)
	if !ok {
		return kit.DefaultAuth{}
	}
	return value
}

// URL is a view helper that returns the current URL.
// The request path can be acccessed with:
// view.URL(ctx).Path
func URL(ctx context.Context) *url.URL {
	value, ok := ctx.Value(middleware.RequestURLKey{}).(*url.URL)
	if !ok {
		return &url.URL{}
	}
	return value

}
