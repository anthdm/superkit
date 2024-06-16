package view

import (
	"context"
	"fmt"
	"net/url"

	"github.com/anthdm/superkit/kit"
	"github.com/anthdm/superkit/kit/middleware"
)

// Asset is view helper that returns the full asset path as a
// string based on the given asset name.
//
//	view.Asset("styles.css") // => /public/assets/styles.css.
func Asset(name string) string {
	return fmt.Sprintf("/public/assets/%s", name)
}

// Auth is a view helper function that returns the current Auth.
// If Auth is not set a default auth will be returned
//
//	view.Auth(ctx)
func Auth(ctx context.Context) kit.Auth {
	value, ok := ctx.Value(kit.AuthKey{}).(kit.Auth)
	if !ok {
		return kit.DefaultAuth{}
	}
	return value
}

// URL is a view helper that returns the current URL.
// The request path can be acccessed with:
//
//	view.URL(ctx).Path // => ex. /login
func URL(ctx context.Context) *url.URL {
	value, ok := ctx.Value(middleware.RequestURLKey{}).(*url.URL)
	if !ok {
		return &url.URL{}
	}
	return value

}
